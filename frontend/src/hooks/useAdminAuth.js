import { useState, useEffect } from 'react';
import axios from 'axios';

const useAdminAuth = () => {

  const [data, setData] = useState({
    loading: true,
    isLoggedIn: false,
  });

  const isAuth = async() => {
    const token = localStorage.getItem('authToken');
    console.log("token", token);

    if (!token) {
      // No token means the user is not logged in
      setData({
        loading: false,
        isLoggedIn: false,
      });
      return;
    }

    await axios.get(
      `${process.env.REACT_APP_BACKEND_URL}/api/v1/admin/is-logged-in`,
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }
    )
    .then(res => {
      if (res.status === 200) {
        setData({
          loading: false,
          isLoggedIn: true,
        });
      }
    })
    .catch(() => {
        setData({
          loading: false,
          isLoggedIn: false,
        });
    });
  }


  useEffect(() => {
    if (data.loading) {
      isAuth();
    }
  },[data]);

  const setLoading = (loading) => {
    setData({
        loading: loading,
        isLoggedIn: data.isLoggedIn,
    })
  }

  return{
    isLoggedIn: data.isLoggedIn,
    token: localStorage.getItem('authToken'),
    loading: data.loading,
    setLoading: setLoading,
  }
}
export default useAdminAuth;