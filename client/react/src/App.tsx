import NavigationBar from "./components/NavBar";
import { useEffect } from "react";
import TaskList from "./components/TaskList";
import Cookies from "js-cookie";

function App() {
  useEffect(() => {
    if (!Cookies.get("token")) {
      window.location.href = "/login";
    }
  });

  return (
    <div className="App">
      <NavigationBar />
      <div className="container mt-4">
        <TaskList />
      </div>
    </div>
  );
}

export default App;
