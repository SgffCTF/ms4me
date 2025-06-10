import { ToastContainer } from 'react-toastify'
import 'bootstrap/dist/css/bootstrap.min.css'
import './styles/App.css'
import { BrowserRouter, Route, Routes } from 'react-router'
import { List } from './pages/list'
import { GameDetail } from './pages/game'
import { Login } from './pages/login'
import { AuthProvider } from './context/AuthProvider'

function App() {
  return (
    <>
      <ToastContainer position="top-right" autoClose={3000} limit={3} />
      <BrowserRouter>
        <AuthProvider>
            <Routes>
                <Route
                path="/"
                element={
                <List />
                }
            />

            <Route
                path="/game/:id"
                element={
                    <GameDetail />
                }
            />
                <Route path="/login" element={<Login />} />
            </Routes>
        </AuthProvider>
      </BrowserRouter>
    </>
  )
}


export default App
