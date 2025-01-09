import React, {Component} from 'react';
import HeaderComponent from "./HeaderComponent";
import FooterComponent from "./FooterComponent";
import {Link} from "react-router-dom";
import img1 from "../assets/images/img1.png"
import img2 from "../assets/images/img2.png"
import img3 from "../assets/images/img3.png"
import axios from "axios";
import authService from "../services/authService";
import RecipesGridComponent from "./RecipesGridComponent";
import "../css/Home.css"

class Home extends Component {
    constructor(props) {
        super(props);
        this.state = {
            recipes: [],
            currentIndexes: {}
        };
    }

    async componentDidMount() {
        await this.fetchRecipes(`http://localhost:8080/api/smart_recipe`);
    }

    fetchRecipes = async (url) => {
        try {
            const accessToken = localStorage.getItem('accessToken');
            const response = await axios.get(url, {
                headers: {
                    Authorization: `Bearer ${accessToken}`
                }
            });
            console.log("Server response:", response);

            if (response.status === 204) {
                this.setState({ hasMoreRecipes: false });
                return;
            }

            if (response.data && Array.isArray(response.data.recipes)) {
                const currentIndexes = {};
                response.data.recipes.forEach(recipe => {
                    currentIndexes[recipe.id] = 0;
                });

                this.setState({
                    recipes: response.data.recipes,
                    currentIndexes,
                });
            } else {
                throw new Error("Invalid response data format");
            }
        } catch (error) {
            console.error("Failed to fetch recipes:", error);
            if (error.response && error.response.status === 401) {
                await authService.refreshToken();
                return this.fetchRecipes(url);
            }
        }
    }

    handlePrevImage = (recipeId) => {
        this.setState(prevState => {
            const currentIndexes = { ...prevState.currentIndexes };
            const recipe = this.state.recipes.find(r => r.id === recipeId);
            if (currentIndexes[recipeId] > 0) {
                currentIndexes[recipeId] -= 1;
            } else {
                currentIndexes[recipeId] = recipe.images.length - 1;
            }
            return { currentIndexes };
        });
    }

    handleNextImage = (recipeId) => {
        this.setState(prevState => {
            const currentIndexes = { ...prevState.currentIndexes };
            const recipe = this.state.recipes.find(r => r.id === recipeId);
            if (currentIndexes[recipeId] < recipe.images.length - 1) {
                currentIndexes[recipeId] += 1;
            } else {
                currentIndexes[recipeId] = 0;
            }
            return { currentIndexes };
        });
    }

    render() {
        return (
            <div className={"home"}>
                <HeaderComponent />
                <div className="hero-container">
                    <div className="hero-content">
                        <h1>Вкусные рецепты &</h1>
                        <h1>Выгодные покупки</h1>
                        <p>Выбирайте, чем себя побаловать, а мы поможем сделать с наибольшей выгодой и удобством</p>
                        <div className={"hero-buttons"}>
                            <Link to={`/recipes`} className={"hero-recipes-button"}>Рецепты</Link>
                        </div>
                    </div>
                    <div className="hero-images">
                        <img className={"hero-image"} alt={"Изображение 1"} src={img1}/>
                        <img className={"hero-image"} alt={"Изображение 2"} src={img2}/>
                        <img className={"hero-image"} alt={"Изображение 3"} src={img3}/>
                    </div>
                </div>
                <h1 className={"main-page-head-lines"}>Недавно добавленное</h1>
                <RecipesGridComponent
                    recipes={this.state.recipes}
                    currentIndexes={this.state.currentIndexes}
                    handlePrevImage={this.handlePrevImage}
                    handleNextImage={this.handleNextImage} />
                <FooterComponent />
            </div>
        );
    }
}

export default Home;
