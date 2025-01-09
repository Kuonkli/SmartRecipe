import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import "./main.css"
import Auth from "./components/Auth";
import Home from './components/Home';
import Recipes from "./components/Recipes";
import Profile from "./components/Profile";
import Recipe from "./components/Recipe";
import NearestShop from "./components/NearestShop";

function App() {
    return (
        <Router>
            <Routes>
                <Route exact path="/" element={<Auth />} />
                <Route exact path="/smart_recipe" element={<Home />} />
                <Route exact path="/recipes" element={<Recipes />} />
                <Route exact path="/recipe/:id" element={<Recipe />} />
                <Route exact path="/recipe/:id/shop" element={<NearestShop />} />
                <Route exact path="/profile" element={<Profile />} />
            </Routes>
        </Router>
    );
}

export default App;
