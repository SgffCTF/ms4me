import { ToastContainer } from 'react-toastify'
import 'bootstrap/dist/css/bootstrap.min.css'
import './App.css'
import { BrowserRouter, Route, Routes } from 'react-router'
import { Board } from './pages/board'
import { Game } from './pages/game'
import { Login } from './pages/login'
import { AuthProvider } from './context/AuthProvider'
import { WebSocketProvider } from './context/WebsocketProvider'

function App() {
  return (
    <>
      <ToastContainer position="top-right" autoClose={3000} />
      <BrowserRouter>
        <AuthProvider>
          <WebSocketProvider>
            <Routes>
              <Route path="/" element={<Board />} />
              <Route path="/game/:id" element={<Game />} />
              <Route path="/login" element={<Login />} />
            </Routes>
          </WebSocketProvider>
        </AuthProvider>
      </BrowserRouter>
    </>
  )
}


export default App
