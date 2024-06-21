import { useState, useEffect } from "react";
import axios from "axios";
import { Table } from "react-bootstrap";

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
        let tasks: Task[] = response.data;
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
      <Table striped bordered hover>
        <thead>
          <tr>
            <th>Task Name</th>
            <th>Last Name</th>
            <th>Username</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>1</td>
            <td>Mark</td>
            <td>Otto</td>
            <td>@mdo</td>
          </tr>
          <tr>
            <td>2</td>
            <td>Jacob</td>
            <td>Thornton</td>
            <td>@fat</td>
          </tr>
          <tr>
            <td>3</td>
            <td colSpan={2}>Larry the Bird</td>
            <td>@twitter</td>
          </tr>
        </tbody>
      </Table>
    </div>
  );
}

export default TaskList;
