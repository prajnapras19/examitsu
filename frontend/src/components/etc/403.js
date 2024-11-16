import React from 'react';
import { Container, Row, Col, Button } from 'react-bootstrap';
import { FaBan } from 'react-icons/fa';

const ForbiddenPage = () => {
  return (
    <Container className="text-center mt-5">
      <Row>
        <Col>
          <h1 className="display-1 text-danger">
            <FaBan /> 403
          </h1>
          <h2 className="mt-3">Tidak ada akses</h2>
          <p className="lead">
            Mohon maaf, Anda tidak memiliki akses ke halaman ini.
          </p>
        </Col>
      </Row>
    </Container>
  );
};

export default ForbiddenPage;
