import { ToastContainer } from 'react-toastify'
import 'bootstrap/dist/css/bootstrap.min.css'
import './styles/App.css'
import { BrowserRouter, Route, Routes } from 'react-router'
import { List } from './pages/list'
import { GameDetail } from './pages/game'
import { Login } from './pages/login'
import { AuthProvider } from './context/AuthProvider'
import { ListWSProvider } from './context/ListWSProvider'
import { GameWSProvider } from './context/GameWSProvider'

function App() {
  return (
    <>
      <ToastContainer position="top-right" autoClose={3000} />
      <BrowserRouter>
        <AuthProvider>
            <Routes>
              <Route
              path="/"
              element={
                <ListWSProvider>
                  <List />
                </ListWSProvider>
              }
            />

            <Route
              path="/game/:id"
              element={
                <GameWSProvider>
                  <GameDetail />
                </GameWSProvider>
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
