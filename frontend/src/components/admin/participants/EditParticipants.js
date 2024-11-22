import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import { toast } from 'react-toastify';
import { Form, Container, Spinner, Button } from 'react-bootstrap';
import ReadExamsMenuCard from '../exam/ReadExamsMenuCard';
import BackToHomepageCard from '../home/BackToHomepageCard';
import ReadParticipantsOfThisExamMenuCard from './ReadParticipantsOfThisExamMenuCard';
import ReadQuestionCard from '../question/ReadQuestionCard';

const EditParticipant = (props) => {
  const { auth } = props;
  const navigate = useNavigate();
  const { examSerial, participantId } = useParams();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [fetchedParticipant, setFetchedParticipant] = useState({});

  const fields = [
    {
      label: 'Kode Peserta',
      name: 'name',
      type: 'text',
      required: true,
      defaultValue: ''
    },
  ]

  // a copy of the fields default value
  const defaultValueMap = {
    name: '',
  }

  const [formData, setFormData] = useState(
    fields.reduce((acc, field) => ({ ...acc, [field.name]: field.defaultValue }), {})
  );

  const handleInputChange = (e) => {
    const { name, type } = e.target;
    if (type === 'checkbox') {
      const { checked } = e.target;
      setFormData({ ...formData, [name]: checked });
    } else {
      const { value } = e.target;
      setFormData({ ...formData, [name]: value });
    }
  };

  useEffect(() => {
    setFormData(
      fields.reduce((acc, field) => ({ ...acc, [field.name]: fetchedParticipant[field.name] ? fetchedParticipant[field.name]: defaultValueMap[field.name]}), {})
    );
    // eslint-disable-next-line
  }, [fetchedParticipant]);

  useEffect(() => {
    if (auth.loading) {
      return;
    }
    if (!auth.isLoggedIn) {
      navigate('/login');
    }

    axios.post(`${process.env.REACT_APP_BACKEND_URL}/api/v1/admin/participants/id/${participantId}`, {
    }, {
      headers: {
        'Authorization': `Bearer ${auth.token}`
      },
    })
    .then(response => {
      if (response.data.data) {
        setFetchedParticipant(response.data.data);
      } else {
        navigate('/404');
      }
      setLoading(false);
    })
    .catch(error => {
      setError(error.message);
      setLoading(false);
    });
  }, [auth.loading, auth.isLoggedIn, auth.token]);

  if (auth.loading || loading) {
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

  const handleSubmit = async (e) => {
    e.preventDefault();

    const customObject = {
      ...formData,
    };

    try {
      await axios.patch(`${process.env.REACT_APP_BACKEND_URL}/api/v1/admin/participants/${participantId}`, customObject, {
        headers: {
          Authorization: `Bearer ${auth.token}`,
        },
      });

      toast.success('Peserta berhasil diubah!', {
        position: "top-center",
        autoClose: 3000,
        hideProgressBar: false,
        closeOnClick: true,
        pauseOnHover: true,
        draggable: true,
      });
      navigate(`/admin/exams/${examSerial}/participants`);
    } catch (err) {
      toast.error(`Gagal mengubah peserta, silakan coba menggunakan data yang lain.`, {
        position: "top-center",
        autoClose: 5000,
        hideProgressBar: false,
        closeOnClick: true,
        pauseOnHover: true,
        draggable: true,
      });
    }
  };

  return (
    <Container>
      <h1 className="my-4">Ubah Peserta</h1>
      <hr/>
      <Container className="text-center mt-5">
        <Container className="card-grid">
          <BackToHomepageCard></BackToHomepageCard>
          <ReadExamsMenuCard></ReadExamsMenuCard>
          <ReadQuestionCard></ReadQuestionCard>
          <ReadParticipantsOfThisExamMenuCard></ReadParticipantsOfThisExamMenuCard>
        </Container>
      </Container>
      <hr/>

      <Form className="my-4" onSubmit={handleSubmit}>
        {fields.map((field) => (
          <Form.Group className="my-3" controlId={field.name} key={field.name}>
            <Form.Label><b>{field.label}{field.required ? <span> (harus diisi)</span> : ''}</b></Form.Label>
            {
              field.type === 'float'
              ? (
                <Form.Control
                  type='number'
                  min={field.minValue}
                  step={field.step}
                  name={field.name}
                  value={formData[field.name]}
                  onChange={handleInputChange}
                  required={field.required}
                  placeholder={`Masukkan ${field.label}`}
                  autoComplete='off'
                />
              )
              : field.type === 'boolean'
              ? (
                <Form.Check
                  type="switch"
                  name={field.name}
                  checked={formData[field.name]}
                  onChange={handleInputChange}
                />
              )
              : (
                <Form.Control
                  type={field.type}
                  name={field.name}
                  value={formData[field.name]}
                  onChange={handleInputChange}
                  required={field.required}
                  placeholder={`Masukkan ${field.label}`}
                  autoComplete='off'
                />
              )
            }
          </Form.Group>
        ))}
        <Button variant="primary" type="submit" className="mt-3">
          Kirim
        </Button>
      </Form>

    </Container>
  );
}

export default EditParticipant;