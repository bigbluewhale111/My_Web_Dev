import { useEffect, useState } from "react";
import axios from "axios";
import Cookies from "js-cookie";

function AppRedirect() {
  const [authenticated, setAuthenticated] = useState(false);
  const [logedout, setLogedout] = useState(false);
  const handlingRedirect = (OauthClientURL: string) => {
    const urlParams = new URLSearchParams(window.location.search);
    var logout = urlParams.get("logout");
    if (logout) {
      axios
        .get(OauthClientURL + "/logout?session=" + Cookies.get("token"))
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
    var code = urlParams.get("code");
    if (code) {
      axios
        .get(OauthClientURL + "/callback?code=" + code)
        .then((response) => {
          var CookieToken = response.data;
          Cookies.set("token", CookieToken);
          setAuthenticated(true);
        })
        .catch((error) => {
          console.error("Error fetching token:", error);
        });
    }
  };
  useEffect(() => {
    axios
      .get("/auth/getOauthClientURL")
      .then((response) => {
        console.log("OauthClientURL:", response.data);
        handlingRedirect(response.data);
      })
      .catch((error) => {
        console.error("Error fetching OauthClientURL:", error);
        return;
      });
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
