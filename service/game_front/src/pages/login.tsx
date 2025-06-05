import { useEffect, useRef, useState } from 'react';
import { fetchLogin, fetchRegister, fetchUser } from '../api/user';
import { toast } from 'react-toastify';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthProvider';

export const Login = () => {
    const [mode, setMode] = useState('login'); // 'login' или 'register'
    const { setUser, user } = useAuth();
    const navigate = useNavigate();
    const username = useRef<HTMLInputElement>(null);
    const password = useRef<HTMLInputElement>(null);

    useEffect(() => {
        if (user != null) {
            navigate("/");
        }
    }, []);

    const handleSubmit = async (e: React.FormEvent<HTMLButtonElement>) => {
        e.preventDefault();
        if (!username || !password) return;
        if (!username.current || !password.current) return;
        
        try {
            if (mode == "login") {
                try {
                    await fetchLogin(username.current.value, password.current.value);
                } catch (err) {
                    if (err instanceof Error) toast.error(err.message); return;
                }
                toast("Успешный вход");
                const user = await fetchUser()
                setUser(user);
                navigate("/");
            } else if (mode == "register") {
                try {
                    await fetchRegister(username.current.value, password.current.value);
                } catch (err) {
                    if (err instanceof Error) toast.error(err.message); return;
                }
                toast("Успешная регистрация");
                setMode("login");
            }
        } catch (e: any) {
            toast.error(e.message);
        }
    }

    return (
        <div className="container d-flex align-items-center justify-content-center h-all">
        <div className="form-signin">
            <div className="d-flex justify-content-center mb-4">
            <button
                className={`btn me-2 ${mode === 'login' ? 'btn-primary' : 'btn-outline-primary'}`}
                onClick={() => setMode('login')}
            >
                Вход
            </button>
            <button
                className={`btn ${mode === 'register' ? 'btn-primary' : 'btn-outline-primary'}`}
                onClick={() => setMode('register')}
            >
                Регистрация
            </button>
            </div>

            <form className="form-signin">
            <h1 className="h3 mb-3 fw-normal text-center">
                {mode === 'login' ? 'Вход' : 'Регистрация'}
            </h1>

            <div className="form-floating mb-3">
                <input type="username" className="form-control" id="floatingInput" placeholder="Имя пользователя" ref={username}/>
                <label htmlFor="floatingInput">Имя пользователя</label>
            </div>

            <div className="form-floating mb-3">
                <input type="password" className="form-control" id="floatingPassword" placeholder="Пароль" ref={password}/>
                <label htmlFor="floatingPassword">Пароль</label>
            </div>

            {/* {mode === 'register' && (
                <div className="form-floating mb-3">
                <input type="password" className="form-control" id="floatingPasswordConfirm" placeholder="Повторите пароль" />
                <label htmlFor="floatingPasswordConfirm">Повторите пароль</label>
                </div>
            )} */}

            <button className="btn btn-primary w-100 py-2" type="submit" onClick={handleSubmit}>
                {mode === 'login' ? 'Войти' : 'Зарегистрироваться'}
            </button>
            </form>
        </div>
        </div>
    );
}
