import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import { toast } from 'react-toastify';
import { Container, Spinner, Table, Button } from 'react-bootstrap';
import BackToHomepageCard from '../home/BackToHomepageCard';
import AddExamCard from './AddExamCard';
import DeleteConfirmationModal from '../../etc/DeleteConfirmationModal';
import DownloadExamTemplateCard from './DownloadExamTemplateCard';
import UploadExamFileCard from './UploadExamFileCard';

const ReadExams = (props) => {
  const { auth } = props;
  const navigate = useNavigate();

  const [currentPage, setCurrentPage] = useState(1);
  const [data, setData] = useState([]);
  const [error, setError] = useState(null);
  const [triggerRender, setTriggerRender] = useState(false);

  const [loadingUpload, setLoadingUpload] = useState(false);
  const onFinishUpload = () => {
    setLoadingUpload(false);
    setTriggerRender(!triggerRender);
  }
  const onClickUpload = () => {
    setLoadingUpload(true);
  }
  
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [deletedExamSerial, setDeletedExamSerial] = useState("");
  const handleShowDeleteModal = (examSerial) => {
    setDeletedExamSerial(examSerial);
    setShowDeleteModal(true);
  }
  const handleCloseDeleteModal = () => {
    setDeletedExamSerial(0);
    setShowDeleteModal(false);
  }
  const handleDelete = () => {
    axios.delete(`${process.env.REACT_APP_BACKEND_URL}/api/v1/admin/exams/${deletedExamSerial}`, { 
      headers: {
        'Authorization': `Bearer ${auth.token}`
      },
    })
    .then(response => {
      toast.success('Ujian berhasil dihapus!', {
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
      toast.error(`Gagal menghapus ujian. Silakan coba beberapa saat lagi.`, {
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
  }, [auth.loading, auth.isLoggedIn]);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.post(`${process.env.REACT_APP_BACKEND_URL}/api/v1/admin/exams?page=${currentPage}`,
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
  }, [currentPage, auth.token, triggerRender]);

  if (auth.loading || loadingUpload) {
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

  const handleNextPage = () => {
    setCurrentPage(currentPage + 1);
  };

  const handlePreviousPage = () => {
    if (currentPage > 1) {
      setCurrentPage(currentPage - 1);
    }
  };

  const copyExamLinkToClipboard = (examSerial) => {
    navigator.clipboard.writeText(`${process.env.REACT_APP_HOST_BASE_URL}/exam/${examSerial}`);
    toast.success('Link ujian berhasil disalin!', {
      position: "top-center",
      autoClose: 3000,
      hideProgressBar: false,
      closeOnClick: true,
      pauseOnHover: true,
      draggable: true,
    });
  }

  return (
    <Container>
      <h1 className="my-4">Daftar Ujian</h1>
      <hr/>
      <Container className="text-center mt-5">
        <Container className="card-grid">
          <BackToHomepageCard></BackToHomepageCard>
          <AddExamCard></AddExamCard>
          <DownloadExamTemplateCard auth={auth}></DownloadExamTemplateCard>
          <UploadExamFileCard
            auth={auth}
            onClick={onClickUpload}
            onFinish={onFinishUpload}
          ></UploadExamFileCard>
        </Container>
      </Container>
      <hr/>


      {data.length === 0 ? (
        <Container className="text-center mt-5">
          {currentPage === 1 ? (
            <i>Tidak ada data ditemukan.</i>  
          ): (
            <i>Tidak ada data ditemukan pada halaman ini. Silakan coba halaman sebelumnya.</i>
          )}
        </Container>
      ) : (
        <Table striped bordered hover className="text-center mt-5">
          <thead>
            <tr>
              <th>Serial</th>
              <th>Nama</th>
              <th>Durasi Pengerjaan (menit)</th>
              <th>Sudah / Masih Bisa Dikerjakan?</th>
              <th colSpan="5">Aksi</th>
            </tr>
          </thead>
          <tbody>
            {data.map((exam) => (
              <tr key={exam.serial}>
                <td className="p-3">{exam.serial}</td>
                <td className="p-3">{exam.name}</td>
                <td className="p-3">{exam.allowed_duration_minutes}</td>
                <td className="p-3">{exam.is_open ? "Ya" : "Tidak"}</td>
                <td>
                  <Button variant="primary" className="me-3" onClick={() => navigate(`/admin/exams/${exam.serial}/edit`)}>Ubah</Button>
                </td>
                <td>
                  <Button variant="primary" className="me-3" onClick={() => copyExamLinkToClipboard(exam.serial)}>Salin Link Ujian</Button>
                </td>
                <td>
                  <Button variant="secondary" className="me-3" onClick={() => navigate(`/admin/exams/${exam.serial}/questions`)}>Lihat Daftar Soal</Button>
                </td>
                <td>
                  <Button variant="secondary" className="me-3" onClick={() => navigate(`/admin/exams/${exam.serial}/participants`)}>Lihat Daftar Peserta</Button>
                </td>
                <td>
                  <Button variant="danger" onClick={() => handleShowDeleteModal(exam.serial)}>Hapus</Button>
                </td>
              </tr>
            ))}
          </tbody>
        </Table>
      )}
      <Container className="d-flex mt-3">
        <Button variant="primary" onClick={handlePreviousPage} disabled={currentPage === 1} className="me-3">
          Halaman sebelumnya
        </Button>
        <Button variant="primary" onClick={handleNextPage} disabled={data.length === 0}>
          Halaman berikutnya
        </Button>
      </Container>
      
      <DeleteConfirmationModal
        show={showDeleteModal}
        handleClose={handleCloseDeleteModal}
        handleDelete={handleDelete}
      />
    </Container>
  );
}

export default ReadExams;