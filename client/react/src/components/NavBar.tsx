import { Navbar, Nav, Container } from "react-bootstrap";

function NavigationBar() {
  return (
    <Navbar bg="dark" variant="dark" expand="lg">
      <Container>
        <Navbar.Brand href="/">
          <img
            src="/checklist.png"
            width="30"
            height="30"
            className="d-inline-block align-top"
            alt="Task logo"
          />
        </Navbar.Brand>
        <Navbar.Toggle aria-controls="basic-navbar-nav" />
        <Navbar.Collapse id="basic-navbar-nav">
          <Nav className="me-auto">
            <Nav.Link href="/">Home</Nav.Link>
            <Nav.Link href="/#">About</Nav.Link>
            <Nav.Link href="/#">Services</Nav.Link>
            <Nav.Link href="/#">Contact</Nav.Link>
          </Nav>
        </Navbar.Collapse>
      </Container>
    </Navbar>
  );
}

export default NavigationBar;
