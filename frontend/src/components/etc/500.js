import React from 'react';
import { Container, Row, Col, Button } from 'react-bootstrap';
import { MdError } from 'react-icons/md';

const InternalServerErrorPage = () => {
  return (
    <Container className="text-center mt-5">
      <Row>
        <Col>
          <h1 className="display-1 text-danger">
            <MdError /> 500
          </h1>
          <h2 className="mt-3">Server Bermasalah</h2>
          <p className="lead">
            Mohon maaf, sedang terjadi masalah pada server. Silakan coba beberapa saat lagi.
          </p>
          <Button variant="primary" href="/">
            Kembali ke halaman utama
          </Button>
        </Col>
      </Row>
    </Container>
  );
};

export default InternalServerErrorPage;
