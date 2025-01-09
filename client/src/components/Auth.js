import React, { Component } from 'react';
import axios from 'axios';
import ("../css/Auth.css")

class Auth extends Component {
    constructor(props) {
        super(props);
        this.state = {
            isLogin: true, // Отслеживаем текущую форму (авторизация или регистрация)
            email: '',
            password: '',
            name: '',
            surname: '',
            errorMessage: ''
        };
        if (localStorage.getItem('accessToken')) {
            window.location.href = "/smart_recipe"
        }
    }

    toggleForm = () => {
        this.setState(prevState => ({
            isLogin: !prevState.isLogin,
            email: '',
            password: '',
            name: '',
            surname: '',
            errorMessage: ''
        }));
    }

    handleChange = (event) => {
        const { name, value } = event.target;
        this.setState({ [name]: value });
    }

    handleSubmit = async (event) => {
        event.preventDefault();
        const { isLogin, email, password, name, surname } = this.state;

        try {
            const response = isLogin
                ? await axios.post('http://localhost:8080/login', { email, password })
                : await axios.post('http://localhost:8080/signup', { name, surname, email, password });

            // Сохранение токенов в локальное хранилище
            const accessToken = response.headers['authorization']?.split(' ')[1];
            const refreshToken = response.headers['x-refresh-token'];

            if (accessToken && refreshToken) {
                localStorage.setItem('accessToken', accessToken);
                localStorage.setItem('refreshToken', refreshToken);

                // Перенаправление после успешной авторизации/регистрации
                window.location.href = '/smart_recipe';
            } else {
                throw new Error('Токены не получены');
            }
        } catch (error) {
            const errorMessage = error.response?.data?.error || 'Произошла ошибка';
            this.setState({ errorMessage });
        }
    }

    render() {
        const { isLogin, email, password, name, surname, errorMessage } = this.state;

        return (
            <div className="authorization">
                <div className="auth-container">
                    <div className="form">
                        <h2>{isLogin ? 'Авторизация' : 'Регистрация'}</h2>
                        {errorMessage && <p className="error-message">{errorMessage}</p>}
                        <form onSubmit={this.handleSubmit}>
                            {!isLogin && (
                                <>
                                    <div className="form-group">
                                        <label htmlFor="name">Имя:</label>
                                        <input
                                            type="text"
                                            id="name"
                                            name="name"
                                            value={name}
                                            onChange={this.handleChange}
                                            required
                                        />
                                    </div>
                                    <div className="form-group">
                                        <label htmlFor="surname">Фамилия:</label>
                                        <input
                                            type="text"
                                            id="surname"
                                            name="surname"
                                            value={surname}
                                            onChange={this.handleChange}
                                            required
                                        />
                                    </div>
                                </>
                            )}
                            <div className="form-group">
                                <label htmlFor="email">Email:</label>
                                <input
                                    type="email"
                                    id="email"
                                    name="email"
                                    value={email}
                                    onChange={this.handleChange}
                                    required
                                />
                            </div>
                            <div className="form-group">
                                <label htmlFor="password">Пароль:</label>
                                <input
                                    type="password"
                                    id="password"
                                    name="password"
                                    value={password}
                                    onChange={this.handleChange}
                                    required
                                />
                            </div>
                            <button className={"auth-button"}
                                    type="submit">{isLogin ? 'Войти' : 'Зарегистрироваться'}</button>
                        </form>
                        <p onClick={this.toggleForm} className="toggle-link">
                            {isLogin ? 'Нет аккаунта? Зарегистрируйтесь' : 'Уже есть аккаунт? Войдите'}
                        </p>
                    </div>
                </div>
            </div>
        );
    }
}

export default Auth;
