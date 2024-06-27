import { useEffect, useState } from "react";
import axios from "axios";
import Cookies from "js-cookie";

function AppRedirect() {
  const [authenticated, setAuthenticated] = useState(false);
  const [logedout, setLogedout] = useState(false);
  useEffect(() => {
    const urlParams = new URLSearchParams(window.location.search);
    var logout = urlParams.get("logout");
    if (logout) {
      axios
        .get("/api/logout")
        .then((response) => {
          console.log(response.data);
          Cookies.remove("token");
          setLogedout(true);
        })
        .catch((error) => {
          console.error("Error logging out:", error);
        })
        .finally(() => {
          return;
        });
    }
    var token = urlParams.get("token");
    if (token) {
      axios
        .get("/auth/callback?token=" + token)
        .then((response) => {
          var CookieToken = response.data;
          Cookies.set("token", CookieToken);
          setAuthenticated(true);
        })
        .catch((error) => {
          console.error("Error fetching token:", error);
        });
    }
  }, []);
  return (
    <>
      {logedout ? (
        <>
          <p>Logged out, redirect soon ...</p>
          <meta http-equiv="refresh" content="2; url=/login" />
        </>
      ) : authenticated ? (
        <>
          <p>Authenticated, redirect soon ...</p>
          <meta http-equiv="refresh" content="2; url=/" />
        </>
      ) : (
        <p>Processing ...</p>
      )}
    </>
  );
}

export default AppRedirect;
