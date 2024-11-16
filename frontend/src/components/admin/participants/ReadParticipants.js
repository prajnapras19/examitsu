import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useParams } from 'react-router-dom';
import { useNavigate } from 'react-router-dom';
import { toast } from 'react-toastify';
import { Container, Spinner, Table, Button } from 'react-bootstrap';
import BackToHomepageCard from '../home/BackToHomepageCard';
import DeleteConfirmationModal from '../../etc/DeleteConfirmationModal';
import ReadExamsMenuCard from '../exam/ReadExamsMenuCard';

const ReadParticipants = (props) => {
  const { auth } = props;
  const { examSerial } = useParams();
  const navigate = useNavigate();

  const [exam, setExam] = useState({});
  const [data, setData] = useState([]);
  const [error, setError] = useState(null);
  const [triggerRender, setTriggerRender] = useState(false);

  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [deletedQuestionId, setDeletedQuestionId] = useState(0);
  const handleShowDeleteModal = (participantId) => {
    setDeletedQuestionId(participantId);
    setShowDeleteModal(true);
  }
  const handleCloseDeleteModal = () => {
    setDeletedQuestionId(0);
    setShowDeleteModal(false);
  }
  const handleDelete = () => {
    axios.delete(`${process.env.REACT_APP_BACKEND_URL}/api/v1/admin/participants/${deletedQuestionId}`, { 
      headers: {
        'Authorization': `Bearer ${auth.token}`
      },
    })
    .then(response => {
      toast.success('Soal berhasil dihapus!', {
        position: "top-center",
        autoClose: 3000,
        hideProgressBar: false,
        closeOnClick: true,
        pauseOnHover: true,
        draggable: true,
      });
      setTriggerRender(!triggerRender);
      handleCloseDeleteModal();
    })
    .catch(err => {
      toast.error(`Gagal menghapus peserta. Silakan coba beberapa saat lagi.`, {
        position: "top-center",
        autoClose: 5000,
        hideProgressBar: false,
        closeOnClick: true,
        pauseOnHover: true,
        draggable: true,
      });
      handleCloseDeleteModal();
    });
  }

  useEffect(() => {
    if (auth.loading) {
      return;
    }
    if (!auth.isLoggedIn) {
      navigate('/admin/login');
    }

    const fetchExam = async() => {
      try {
        const response = await axios.post(`${process.env.REACT_APP_BACKEND_URL}/api/v1/admin/exams/${examSerial}`,
          {}, {
            headers: {
              Authorization: `Bearer ${auth.token}`,
            },
          },
        );

        setExam(response.data.data);
      } catch (err) {
        console.error("Error fetching exam", err);
        setError(err);
      }
    }
    
    fetchExam();
  }, [auth.loading, auth.isLoggedIn]);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.post(`${process.env.REACT_APP_BACKEND_URL}/api/v1/admin/participants/exam-serial/${examSerial}`,
          {}, {
            headers: {
              Authorization: `Bearer ${auth.token}`,
            },
          },
        );

        setData(response.data.data);
      } catch (err) {
        console.error("Error fetching data", err);
        setError(err);
      }
    };

    fetchData();
  }, [auth.token, triggerRender]);

  if (auth.loading) {
    return (
      <Container className="text-center">
        <Spinner animation="border" />
        <p>Mohon tunggu...</p>
      </Container>
    );
  }

  if (error) {
    navigate('/500');
  }

  return (
    <Container>
      <h1 className="my-4">Daftar Peserta</h1>
      <hr/>
      <Container className="text-center mt-5">
        <Container className="card-grid">
          <BackToHomepageCard></BackToHomepageCard>
          <ReadExamsMenuCard></ReadExamsMenuCard>
        </Container>
      </Container>
      <hr/>

      <Container className="mt-5">
        <h3>Ujian</h3>
        <Table striped bordered hover className="text-center mt-5">
          <thead>
            <tr>
              <th>Serial</th>
              <th>Nama</th>
              <th>Sudah / Masih Bisa Dikerjakan?</th>
            </tr>
          </thead>
          <tbody>
            <tr key={exam.serial}>
              <td>{exam.serial}</td>
              <td>{exam.name}</td>
              <td>{exam.is_open ? "Ya" : "Tidak"}</td>
            </tr>
          </tbody>
        </Table>
      </Container>
      <hr/>

      <h3>Peserta</h3>
      {data.length === 0 ? (
        <Container className="text-center mt-5">
          <i>Tidak ada data ditemukan.</i>  
        </Container>
      ) : (
        <Table striped bordered hover className="text-center mt-5">
          <thead>
            <tr>
              <th>#</th>
              <th>Nama</th>
              <th>Kata Sandi</th>
              <th colSpan="3">Aksi</th>
            </tr>
          </thead>
          <tbody>
            {data.map((participant, i) => (
              <tr key={participant.id}>
                <td className="p-3">{i+1}</td>
                <td className="p-3">{participant.name}</td>
                <td className="p-3">{participant.password}</td>
                <td>
                <Button variant="primary" className="me-3" onClick={() => navigate(`/admin/exams/${exam.serial}/participants/${participant.id}/edit`)}>Ubah</Button>
                </td>
                <td>
                  <Button variant="danger" onClick={() => handleShowDeleteModal(participant.id)}>Hapus</Button>
                </td>
              </tr>
            ))}
          </tbody>
        </Table>
      )}
      
      <DeleteConfirmationModal
        show={showDeleteModal}
        handleClose={handleCloseDeleteModal}
        handleDelete={handleDelete}
      />
    </Container>
  );
}

export default ReadParticipants;