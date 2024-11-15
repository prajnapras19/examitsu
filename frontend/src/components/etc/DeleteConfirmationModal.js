import React from 'react';
import { Button, Modal } from 'react-bootstrap';

const DeleteConfirmationModal = ({ show, handleClose, handleDelete }) => {
  return (
    <Modal show={show} onHide={handleClose}>
      <Modal.Header closeButton>
        <Modal.Title>Konfirmasi Penghapusan</Modal.Title>
      </Modal.Header>
      <Modal.Body>Apakah Anda yakin ingin menghapus data ini?</Modal.Body>
      <Modal.Footer>
        <Button variant="secondary" onClick={handleClose}>
          Tidak
        </Button>
        <Button variant="danger" onClick={handleDelete}>
          Iya
        </Button>
      </Modal.Footer>
    </Modal>
  );
};

export default DeleteConfirmationModal;
