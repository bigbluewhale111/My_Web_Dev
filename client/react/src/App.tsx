import NavigationBar from "./components/NavBar";
import TaskList from "./components/TaskList";
function App() {
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
