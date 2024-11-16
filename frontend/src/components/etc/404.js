import React from 'react';
import { Container, Row, Col, Button } from 'react-bootstrap';
import { TbWorldQuestion } from "react-icons/tb";

const NotFoundPage = () => {
  return (
    <Container className="text-center mt-5">
      <Row>
        <Col>
          <h1 className="display-1 text-danger">
            <TbWorldQuestion /> 404
          </h1>
          <h2 className="mt-3">Halaman tidak ditemukan</h2>
          <p className="lead">
            Halaman yang ingin Anda tuju tidak ditemukan.
          </p>
        </Col>
      </Row>
    </Container>
  );
};

export default NotFoundPage;
