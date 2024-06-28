import React from 'react';
import './App.css';
import { Navbar } from './components/Navbar';
import { Route, BrowserRouter as Router, Routes } from 'react-router-dom';
import Clients from './pages/Clients';
import Home from './pages/Home';
import Client from './pages/Client';
import Vehicles from './pages/Vehicles';
import Error404 from './pages/404';

function App() {
  return (
    <div className="App">
      <Router>
        <Navbar />
        <Routes>
          <Route path='/' element={<Home />} />
          <Route path='/clients' element={<Clients />} />
          <Route path='/vehicles/:id' element={<Vehicles />} />
          <Route path='/clients/:id' element={<Client />} />
          <Route path='*' element={<Error404 />} />
        </Routes>
      </Router>
    </div>
  );
}

export default App;
