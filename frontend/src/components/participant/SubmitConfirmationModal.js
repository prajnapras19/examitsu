import React from 'react';
import { Button, Modal } from 'react-bootstrap';

const SubmitConfirmationModal = ({ show, handleClose, handleSubmit }) => {
  return (
    <Modal show={show} onHide={handleClose}>
      <Modal.Header closeButton>
        <Modal.Title>Konfirmasi Pengumpulan</Modal.Title>
      </Modal.Header>
      <Modal.Body>Apakah Anda yakin ingin mengumpulkan ujian? Anda tidak akan bisa mengubah jawaban setelah dikumpulkan.</Modal.Body>
      <Modal.Footer>
        <Button variant="secondary" onClick={handleClose}>
          Tidak
        </Button>
        <Button variant="danger" onClick={handleSubmit}>
          Iya
        </Button>
      </Modal.Footer>
    </Modal>
  );
};

export default SubmitConfirmationModal;
