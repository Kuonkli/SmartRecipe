import React from 'react';
import { Link } from 'react-router-dom';
import ("../css/FooterComponent.css")
import InstagramIcon from '../assets/images/instagram.png';
import FacebookIcon from '../assets/images/facebook.png';
import TwitterIcon from '../assets/images/twitter.png';

const Footer = () => {
    return (
        <footer>
            <div className="footer-content">
                <div className="footer-logo">
                    <h3>Smart.Recipe</h3>
                    <p>Готовьте вкусно и недорого</p>
                </div>
                <div className="footer-links">
                    <ul>
                        <li><Link to="/">О нас</Link></li>
                        <li><Link to="/">Контакты</Link></li>
                        <li><Link to="/">Политика конфиденциальности</Link></li>
                        <li><Link to="/">Условия использования</Link></li>
                    </ul>
                </div>
                <div className="footer-social">
                    <p>Контакты</p>
                    <div className="social-icons">
                        <Link to="/"><img src={InstagramIcon} alt="Instagram" /></Link>
                        <Link to="/"><img src={FacebookIcon} alt="Facebook" /></Link>
                        <Link to="/"><img src={TwitterIcon} alt="Twitter" /></Link>
                    </div>
                    <p>2024 @smartrecipe.com</p>
                </div>
            </div>
            <p className={"smart-recipe-text-container"}>2024 © SmartRecipe</p>
        </footer>
    );
};

export default Footer;
