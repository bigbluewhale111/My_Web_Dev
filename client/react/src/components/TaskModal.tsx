import axios from "axios";
import { useState } from "react";
import { Button, Modal } from "react-bootstrap";
import { FaPencil } from "react-icons/fa6";

interface Task {
  id: number;
  name: string;
  description: string;
  due_date: number;
  status: string;
}

function TaskModals(props: { task_id: number }) {
  const [task, setTask] = useState<Task>();
  const [showInfo, setShowInfo] = useState(false);
  const handleShowInfo = () => {
    axios
      .get("http://localhost:3000/api/task/" + props.task_id)
      .then((response) => {
        console.log(response.data);
        let task: Task = response.data;
        setTask(task);
      })
      .catch((error) => {
        console.error("Error fetching task:", error);
      });
    setShowInfo(true);
  };
  const handleCloseInfo = () => setShowInfo(false);
  // const [showEdit, setShowEdit] = useState(false);
  // const handleShowEdit = () => setShowEdit(true);
  // const handleCloseEdit = () => setShowEdit(false);
  return (
    <>
      <Button
        variant="primary"
        style={{ marginRight: "10px" }}
        onClick={handleShowInfo}
      >
        <FaPencil style={{ color: "white" }}></FaPencil>
      </Button>

      <Modal show={showInfo} onHide={handleCloseInfo}>
        <Modal.Header closeButton>
          <Modal.Title>
            {task != undefined ? task.name : "No Task Found"}
          </Modal.Title>
        </Modal.Header>
        {task != undefined ? (
          <Modal.Body>
            <p>{task.description}</p>
            <p>Due Date: {new Date(task.due_date).toLocaleString()}</p>
            <p>Status: {task.status}</p>
          </Modal.Body>
        ) : (
          <></>
        )}
        <Modal.Footer>
          <Button variant="secondary" onClick={handleCloseInfo}>
            Edit
          </Button>
          {task != undefined ? (
            <Button variant="primary" onClick={handleCloseInfo}>
              Close
            </Button>
          ) : (
            <></>
          )}
        </Modal.Footer>
      </Modal>
    </>
  );
}

export default TaskModals;
