import { useState, useEffect } from "react";
import axios from "axios";

interface Task {
  id: number;
  name: string;
}

function TaskList() {
  const [tasks, setTasks] = useState<Task[]>([]);
  // const [newTask, setNewTask] = useState("");

  useEffect(() => {
    axios
      .get("http://localhost:3000/api/tasks")
      .then((response) => {
        console.log(response.data);
        let tasks: Task[] = JSON.parse(response.data);
        setTasks(tasks);
      })
      .catch((error) => {
        console.error("Error fetching tasks:", error);
      });
  }, []);

  // const handleAddTask = () => {
  //   // Make a POST request to add a new task
  //   if (newTask.trim() === "") {
  //     return; // Do nothing if the input is empty
  //   }

  //   // Make a POST request to add a new task
  //   axios
  //     .post("http://localhost:3000/api/tasks", { name: newTask })
  //     .then((response) => {
  //       // Update the state with the new list of tasks returned by the server
  //       setTasks(response.data);
  //       setNewTask(""); // Clear the input field
  //     })
  //     .catch((error) => {
  //       console.error("Error adding task:", error);
  //     });
  // };

  return (
    <div>
      <h2>Task List</h2>
      <ul>
        {
          /* {tasks.length > 0 ? (
          tasks.map((task) => <li key={task.id}>{task.name}</li>)
        ) : (
          <li>No tasks available</li>
        )} */
          tasks.toString()
        }
      </ul>
      {/* <input
        type="text"
        value={newTask}
        onChange={(e) => setNewTask(e.target.value)}
        placeholder="New task name"
      />
      <button onClick={handleAddTask}>Add Task</button> */}
    </div>
  );
}

export default TaskList;
