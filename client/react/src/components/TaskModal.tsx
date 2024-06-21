import axios from "axios";
import { useState } from "react";
import {
  Button,
  Modal,
  Container,
  Row,
  Col,
  Form,
  Alert,
} from "react-bootstrap";
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
  const [validDueDate, setValidDueDate] = useState(true);
  const [newName, setNewName] = useState("");
  const [newDescription, setNewDescription] = useState("");
  const [newDueDate, setNewDueDate] = useState("");
  const [newStatus, setNewStatus] = useState("");
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
  const [showEdit, setShowEdit] = useState(false);
  const handleShowEdit = () => {
    setValidDueDate(true);
    setShowEdit(true);
  };
  const handleCloseEdit = () => setShowEdit(false);
  const handleSendEdit = () => {
    var sendDueDate = Date.parse(newDueDate);
    if (newDueDate != "" && isNaN(sendDueDate)) {
      setValidDueDate(false);
      return;
    }
    if (newDueDate == "") {
      sendDueDate = 0;
    }
    axios
      .post("http://localhost:3000/api/edit/task/" + props.task_id, {
        name: newName,
        description: newDescription,
        due_date: sendDueDate,
        status: newStatus,
      })
      .then((response) => {
        console.log(response.data);
      })
      .catch((error) => {
        console.error("Error editing task:", error);
      })
      .finally(() => {
        handleCloseEdit();
      });
  };
  return (
    <>
      <Button
        variant="primary"
        style={{ marginRight: "10px" }}
        onClick={handleShowInfo}
      >
        <FaPencil style={{ color: "white" }}></FaPencil>
      </Button>

      <Modal show={showInfo} onHide={handleCloseInfo} size="lg">
        <Modal.Header closeButton>
          <Modal.Title>
            {task != undefined ? "Edit Task" : "No Task Found"}
          </Modal.Title>
        </Modal.Header>
        {task != undefined ? (
          <Modal.Body>
            <Form>
              <Container>
                <Row>
                  <Col>
                    <p>Due Date: {new Date(task.due_date).toLocaleString()}</p>
                  </Col>
                  <Col>
                    <p>Status: {task.status}</p>
                  </Col>
                </Row>
                <Row>
                  <p>{task.description}</p>
                </Row>
              </Container>
            </Form>
          </Modal.Body>
        ) : (
          <></>
        )}
        <Modal.Footer>
          {task != undefined ? (
            <Button
              variant="secondary"
              onClick={() => {
                handleCloseInfo();
                handleShowEdit();
              }}
            >
              Edit
            </Button>
          ) : (
            <></>
          )}

          <Button variant="primary" onClick={handleCloseInfo}>
            Cancel
          </Button>
        </Modal.Footer>
      </Modal>

      <Modal show={showEdit} onHide={handleCloseEdit} size="lg">
        <Modal.Header closeButton>
          <Modal.Title>
            {task != undefined ? task.name : "No Task Found"}
          </Modal.Title>
        </Modal.Header>
        {task != undefined ? (
          <Modal.Body>
            <Container>
              <Row>
                <Form.Group controlId="formName">
                  <Form.Label>Name:</Form.Label>
                  <Form.Control
                    placeholder={task.name}
                    onChange={(e) => {
                      setNewName(e.target.value);
                    }}
                  ></Form.Control>
                </Form.Group>
              </Row>
              <Row>
                <Col>
                  <Form.Group controlId="formDueDate">
                    <Form.Label>Due Date:</Form.Label>
                    <Form.Control
                      placeholder={new Date(task.due_date).toLocaleString()}
                      onChange={(e) => {
                        setNewDueDate(e.target.value);
                        setValidDueDate(true);
                      }}
                    ></Form.Control>
                  </Form.Group>
                </Col>
                <Col>
                  <Form.Group controlId="formStatus">
                    <Form.Label>Status:</Form.Label>
                    <Form.Control
                      placeholder={task.status}
                      onChange={(e) => {
                        setNewStatus(e.target.value);
                      }}
                    ></Form.Control>
                  </Form.Group>
                </Col>
              </Row>
              <Row>
                <Form.Group controlId="formDescription">
                  <Form.Control
                    as="textarea"
                    placeholder={task.description}
                    style={{ height: "100px" }}
                    onChange={(e) => {
                      setNewDescription(e.target.value);
                    }}
                  ></Form.Control>
                </Form.Group>
              </Row>
              <Row>
                {validDueDate ? (
                  ""
                ) : (
                  <Alert variant="danger">
                    Please enter a valid date in the format MM/DD/YYYY,HH/MM/SS
                    AM/PM
                  </Alert>
                )}
              </Row>
            </Container>
          </Modal.Body>
        ) : (
          <></>
        )}
        <Modal.Footer>
          {task != undefined ? (
            <Button variant="secondary" onClick={handleSendEdit}>
              Edit
            </Button>
          ) : (
            <></>
          )}
          <Button variant="primary" onClick={handleCloseEdit}>
            Close
          </Button>
        </Modal.Footer>
      </Modal>
    </>
  );
}

export default TaskModals;
