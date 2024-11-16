import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import { toast } from 'react-toastify';
import { Form, Container, Spinner, Button } from 'react-bootstrap';
import ReadExamsMenuCard from './ReadExamsMenuCard';
import BackToHomepageCard from '../home/BackToHomepageCard';

const EditExam = (props) => {
  const { auth } = props;
  const navigate = useNavigate();
  const { examSerial } = useParams();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [fetchedExam, setFetchedExam] = useState({});

  const fields = [
    {
      label: 'Nama Ujian',
      name: 'name',
      type: 'text',
      required: true,
    },
    {
      label: 'Bisa Dikerjakan??',
      name: 'is_open',
      type: 'boolean',
      defaultValue: false,
    },
  ]

  // a copy of the fields default value
  const defaultValueMap = {
    name: '',
    is_open: false,
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
      fields.reduce((acc, field) => ({ ...acc, [field.name]: fetchedExam[field.name] ? fetchedExam[field.name]: defaultValueMap[field.name]}), {})
    );
    // eslint-disable-next-line
  }, [fetchedExam]);

  useEffect(() => {
    if (auth.loading) {
      return;
    }
    if (!auth.isLoggedIn) {
      navigate('/login');
    }

    axios.post(`${process.env.REACT_APP_BACKEND_URL}/api/v1/admin/exams`, {
        serial_equals_to: {
          value: examSerial,
        },
    }, {
      headers: {
        'Authorization': `Bearer ${auth.token}`
      },
    })
    .then(response => {
      if (response.data.data && response.data.data.length === 1) {
        setFetchedExam(response.data.data[0]);
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
      await axios.patch(`${process.env.REACT_APP_BACKEND_URL}/api/v1/admin/exams/${examSerial}`, customObject, {
        headers: {
          Authorization: `Bearer ${auth.token}`,
        },
      });

      toast.success('Ujian berhasil diubah!', {
        position: "top-center",
        autoClose: 3000,
        hideProgressBar: false,
        closeOnClick: true,
        pauseOnHover: true,
        draggable: true,
      });
      navigate('/admin/exams');
    } catch (err) {
      toast.error(`Gagal mengubah ujian, silakan coba menggunakan data yang lain.`, {
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
      <h1 className="my-4">Ubah Ujian</h1>
      <hr/>
      <Container className="text-center mt-5">
        <Container className="card-grid">
          <BackToHomepageCard></BackToHomepageCard>
          <ReadExamsMenuCard></ReadExamsMenuCard>
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

export default EditExam;