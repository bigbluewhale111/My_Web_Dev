import { useEffect, useState } from "react";
import NavigationBar from "./components/NavBar";
import Button from "react-bootstrap/Button";
import { FaGithub } from "react-icons/fa";
import axios from "axios";
// import Cookies from "js-cookie";

function AppLogin() {
  // if (Cookies.get("token")) {
  //   window.location.href = "/";
  // }
  const [clientId, setClientId] = useState("");
  useEffect(() => {
    axios
      .get("/auth/client_id")
      .then((res) => {
        setClientId(res.data);
      })
      .catch((error) => {
        console.error("Error fetching client_id:", error);
      });
  });
  return (
    <div className="AppLogin">
      <NavigationBar />
      <div className="container mt-4"></div>
      <div style={{ textAlign: "center" }}>
        <h1>Login</h1>
        <p>You are not logged in, please login with Github</p>
        {clientId === "" ? (
          <p>Cannot get client_id</p>
        ) : (
          <Button
            variant="success"
            href={
              "https://github.com/login/oauth/authorize?client_id=" + clientId
            }
          >
            Login with Github <FaGithub></FaGithub>
          </Button>
        )}
      </div>
    </div>
  );
}

export default AppLogin;
