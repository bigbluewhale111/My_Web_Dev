import { useEffect, useState } from "react";
import NavigationBar from "./components/NavBar";
import Button from "react-bootstrap/Button";
import { FaGithub } from "react-icons/fa";
import axios from "axios";
import Cookies from "js-cookie";

function AppLogin() {
  const [oauth2Url, setOauth2Url] = useState("");
  useEffect(() => {
    if (Cookies.get("token")) {
      window.location.href = "/";
    }
    axios
      .get("/auth/getOauthURL")
      .then((res) => {
        setOauth2Url(res.data);
      })
      .catch((error) => {
        console.error("Error fetching oauth2 url:", error);
      });
  }, []);
  return (
    <div className="AppLogin">
      <NavigationBar />
      <div className="container mt-4"></div>
      <div style={{ textAlign: "center" }}>
        <h1>Login</h1>
        <p>You are not logged in, please login with Github</p>
        {oauth2Url === "" ? (
          <p>Cannot get Oauth2 URL</p>
        ) : (
          <>
            <Button
              variant="success"
              href={
                oauth2Url +
                "/login?callback_url=" +
                document.location.origin +
                "/redirect"
              }
            >
              Login with My OAUTH <FaGithub></FaGithub>
            </Button>
          </>
        )}
      </div>
    </div>
  );
}

export default AppLogin;
