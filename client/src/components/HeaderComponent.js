import React, { useState, useEffect } from 'react';
import { Link } from "react-router-dom";
import ("../css/HeaderComponent.css")

const HeaderComponent = () => {
    const [isDropdownVisible, setDropdownVisible] = useState(false);
    const [isMobile, setIsMobile] = useState(window.innerWidth <= 768);

    const navLinks = [
        { name: 'Главная', path: '/smart_recipe' },
        { name: 'Рецепты', path: '/recipes' },
        { name: 'Планы питания', path: '/' },
        { name: 'О нас', path: '/' },
        { name: 'Профиль', path: '/profile', className: 'profile-btn' }
    ];

    const toggleDropdown = () => {
        setDropdownVisible(!isDropdownVisible);
    };

    const closeDropdown = () => {
        setDropdownVisible(false);
    };

    const handleResize = () => {
        setIsMobile(window.innerWidth <= 768);
        if (window.innerWidth > 768) {
            setDropdownVisible(false);
        }
    };

    useEffect(() => {
        window.addEventListener('resize', handleResize);
        return () => {
            window.removeEventListener('resize', handleResize);
        };
    }, []);

    return (
        <header>
            <nav className="navbar">
                <a href="/" className="logo">
                    <div className="smart">Smart</div>.Recipe
                </a>
                <ul className="nav-links">
                    {navLinks.map((link, index) => (
                        <li key={index}>
                            <Link to={link.path} className={link.className}>{link.name}</Link>
                        </li>
                    ))}
                </ul>
            </nav>
            {isMobile && (
                <div className="dropdown">
                    <button className="dropdown-button" onClick={toggleDropdown}>Меню</button>
                    <div>
                        {isDropdownVisible && (
                            <ul className="dropdown-menu">
                                {navLinks.map((link, index) => (
                                    <li key={index}>
                                        <Link to={link.path} onClick={closeDropdown}>{link.name}</Link>
                                    </li>
                                ))}
                            </ul>
                        )}
                    </div>
                </div>
            )}
        </header>
    );
};

export default HeaderComponent;
