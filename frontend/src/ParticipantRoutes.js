import { Route, Routes } from 'react-router-dom';
import NotFoundPage from './components/etc/404';

function ParticipantRoutes() {
  document.addEventListener('contextmenu', (e) => {
    e.preventDefault();
  });
  return (
    <Routes>
      <Route path="*" element={<NotFoundPage/>}/>
    </Routes>
  );
}

export default ParticipantRoutes;