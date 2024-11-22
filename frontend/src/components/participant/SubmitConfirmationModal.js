import React from 'react';
import { Button, Modal } from 'react-bootstrap';

const SubmitConfirmationModal = ({ show, handleClose, handleSubmit }) => {
  return (
    <Modal show={show} onHide={handleClose}>
      <Modal.Header closeButton>
        <Modal.Title className="prevent-select">Konfirmasi Pengumpulan</Modal.Title>
      </Modal.Header>
      <Modal.Body className="prevent-select">
        <p>Apakah Anda yakin ingin mengumpulkan ujian?</p>
        <p><b>Pastikan bahwa Anda sudah mengisi semua jawaban!</b></p>
        <p>Anda tidak akan bisa mengubah jawaban setelah dikumpulkan.</p>
      </Modal.Body>
      <Modal.Footer className="prevent-select">
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
