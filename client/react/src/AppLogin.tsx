import NavigationBar from "./components/NavBar";
import Button from "react-bootstrap/Button";
import { FaGithub } from "react-icons/fa";
import Cookies from "js-cookie";

function AppLogin() {
  // if (Cookies.get("token")) {
  //   window.location.href = "/";
  // }
  const handleLogin = () => {
    document.location.assign("https://github.com/auth/github");
  };
  return (
    <div className="AppLogin">
      <NavigationBar />
      <div className="container mt-4"></div>
      <div style={{ textAlign: "center" }}>
        <h1>Login</h1>
        <p>You are not logged in, please login with Github</p>
        <Button variant="success" onClick={handleLogin}>
          Login with Github <FaGithub></FaGithub>
        </Button>
      </div>
    </div>
  );
}

export default AppLogin;
