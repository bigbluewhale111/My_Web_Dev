import { useEffect, useState } from "react";
import axios from "axios";
import Cookies from "js-cookie";

function AppRedirect() {
  const [authenticated, setAthenticated] = useState(false);
  const [logedout, setLogedout] = useState(false);
  useEffect(() => {
    const urlParams = new URLSearchParams(window.location.search);
    var logout = urlParams.get("logout");
    if (logout) {
      axios
        .get("/api/logout")
        .then((response) => {
          setLogedout(true);
          console.log(response.data);
        })
        .catch((error) => {
          console.error("Error logging out:", error);
        });
      Cookies.remove("token");
      return;
    }
    var code = urlParams.get("code");
    axios
      .get("/auth/callback?code=" + code)
      .then((response) => {
        var token = response.data;
        Cookies.set("token", token);
        console.log(token);
        setAthenticated(true);
      })
      .catch((error) => {
        console.error("Error fetching token:", error);
      });
  });
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
