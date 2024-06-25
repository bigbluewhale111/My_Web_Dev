import NavigationBar from "./components/NavBar";
import TaskList from "./components/TaskList";
// import Cookies from "js-cookie";

function App() {
  // if (!Cookies.get("token")) {
  //   window.location.href = "/login";
  // }
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
