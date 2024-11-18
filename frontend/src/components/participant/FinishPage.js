import React from 'react';
import { Container, Row, Col } from 'react-bootstrap';
import { IconContext } from 'react-icons';
import { FaCheckCircle } from "react-icons/fa";

const FinishPage = () => {
  return (
    <Container className="text-center mt-5">
      <Row>
        <Col>
          <h1 className="display-1">
          <IconContext.Provider
            value={{ color: 'green' }}
          >
            <FaCheckCircle />
          </IconContext.Provider>
          </h1>
          <h2 className="mt-3">Ujian berhasil dikumpulkan!</h2>
          <p className="lead">
            Terima kasih sudah mengerjakan. Tetap semangat!
          </p>
        </Col>
      </Row>
    </Container>
  );
};

export default FinishPage;
