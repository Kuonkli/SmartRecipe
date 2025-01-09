import axios from 'axios';

const authService = {
    clearTokens: () => {
        localStorage.removeItem('accessToken');
        localStorage.removeItem('refreshToken');
    },

    async refreshToken() {
        try {
            const refreshToken = localStorage.getItem('refreshToken');
            if (!refreshToken) {
                this.clearTokens()
                window.location.href = "/";
            }
            const response = await axios.post('http://localhost:8080/refresh', {}, {
                headers: {
                    "X-Refresh-Token": refreshToken,
                }
            });
            const newAccessToken = response.headers['authorization'].split(' ')[1];
            localStorage.setItem('accessToken', newAccessToken)
        } catch (error) {
            this.clearTokens();
            window.location.href = "/"
        }
    },
};

export default authService;
