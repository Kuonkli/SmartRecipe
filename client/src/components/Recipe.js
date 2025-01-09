import React, { Component, useEffect } from "react";
import {Link, useParams} from 'react-router-dom';
import "../css/Recipe.css";
import axios from "axios";
import authService from "../services/authService";
import HeaderComponent from "./HeaderComponent";
import FooterComponent from "./FooterComponent";
import FormMarkComponent from "./FormMarkComponent";
import DefaultImage from '../assets/images/default-recipe-image.png';

function withParams(Component) {
    return props => <Component {...props} params={useParams()} />;
}

const ScrollToTop = () => { useEffect(() => {
    window.scrollTo(0, 0); }, []);
    return null;
}

class Recipe extends Component {
    constructor(props) {
        super(props);
        this.state = {
            recipe: {},
            isFavorite: false,
            isTried: false,
            userScore: 0,
            ingredients: [],
            currentIndex: 0,
            instructions: [],
            isFormVisible: false
        };
    }

    async componentDidMount() {
        const { id } = this.props.params;
        await this.fetchRecipe(`http://localhost:8080/api/recipes/${id}`);
    }

    fetchRecipe = async (url) => {
        try {
            const accessToken = localStorage.getItem('accessToken');
            const response = await axios.get(url, {
                headers: {
                    Authorization: `Bearer ${accessToken}`
                }
            });
            console.log(response);
            this.setState({
                recipe: response.data.recipe,
                isFavorite: response.data.is_favorite,
                isTried: response.data.is_tried,
                userScore: response.data.user_score,
                ingredients: response.data.ingredients,
                instructions: this.parseInstructions(response.data.recipe.instructions),
            });
        } catch (error) {
            console.error("Failed to fetch recipe:", error);
            if (error.response && error.response.status === 401) {
                await authService.refreshToken();
                await this.fetchRecipe(url);
            }
        }
    }

    fetchRecipeActions = async (url) => {
        try {
            const accessToken = localStorage.getItem('accessToken');
            const response = await axios.put(url, {
                score: this.state.userScore,
                is_favorite: this.state.isFavorite,
                is_tried: this.state.isTried
            }, {
                headers: {
                    Authorization: `Bearer ${accessToken}`
                }
            });
            console.log(response);
        } catch (error) {
            console.error("Failed to update recipe action:", error);
            if (error.response && error.response.status === 401) {
                await authService.refreshToken();
                await this.fetchRecipeActions(url);
            }
        }
    }

    handleAddFavorite = async () => {
        this.setState(prevState => ({
            isFavorite: !prevState.isFavorite
        }), async () => {
            await this.fetchRecipeActions(`http://localhost:8080/api/recipes/${this.state.recipe.id}/edit/user_recipes`);
        });
    }

    handleAddTried = async () => {
        this.setState(prevState => ({
            isTried: !prevState.isTried
        }), async () => {
            await this.fetchRecipeActions(`http://localhost:8080/api/recipes/${this.state.recipe.id}/edit/user_recipes`);
        });
    }

    handleRatingSubmit = async (rating) => {
        this.setState({
            userScore: rating
        }, async () => {
            await this.fetchRecipeActions(`http://localhost:8080/api/recipes/${this.state.recipe.id}/edit/user_recipes`)
        });
        this.toggleFormVisibility()
    }

    toggleFormVisibility = () => {
        this.setState(prevState => ({
            isFormVisible: !prevState.isFormVisible
        }))
    }

    parseInstructions = (instructions) => {
        if (!instructions) return [];
        return instructions.split(/(\d+\.\s+)/g).filter(step => step.trim() !== "" && !step.match(/^\d+\.\s+$/));
    }

    handlePrevImage = () => {
        this.setState(prevState => ({
            currentIndex: (prevState.currentIndex > 0) ? prevState.currentIndex - 1 : prevState.recipe.images.length - 1
        }));
    }

    handleNextImage = () => {
        this.setState(prevState => ({
            currentIndex: (prevState.currentIndex < prevState.recipe.images.length - 1) ? prevState.currentIndex + 1 : 0
        }));
    }

    handleThumbnailClick = (index) => {
        this.setState({ currentIndex: index });
    }

    getStars(rating) {
        const fullStars = Math.floor(rating);
        const halfStar = rating % 1 >= 0.5 ? "★" : "";
        const emptyStars = "☆".repeat(5 - fullStars - (halfStar ? 1 : 0));
        return "★".repeat(fullStars) + halfStar + emptyStars;
    }

    render() {
        const { recipe, currentIndex, isFavorite, isTried, userScore, ingredients, instructions, isFormVisible } = this.state;
        const hasImages = recipe.images && recipe.images.length > 0;
        const currentImage = hasImages ? `http://localhost:8080/${recipe.images[currentIndex].image_path}` : DefaultImage;
        return (
            <div className="recipe-page">
                <HeaderComponent />
                <ScrollToTop />
                <div className="recipe">
                    <div className="recipe-information-container">
                        <div className="recipe-image-container">
                            <div className="recipe-main-image">
                                <img src={currentImage} alt="Main Recipe Image"/>
                                {hasImages && recipe.images.length > 1 && (
                                    <div className="image-nav-buttons">
                                        <button onClick={this.handlePrevImage}>
                                            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 50">
                                                <line x1="15" y1="10" x2="5" y2="25" stroke="white" strokeWidth="2"
                                                      strokeLinecap="round"/>
                                                <line x1="15" y1="40" x2="5" y2="25" stroke="white" strokeWidth="2"
                                                      strokeLinecap="round"/>
                                            </svg>
                                        </button>
                                        <button onClick={this.handleNextImage}>
                                            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 50">
                                                <line x1="5" y1="10" x2="15" y2="25" stroke="white" strokeWidth="2"
                                                      strokeLinecap="round"/>
                                                <line x1="5" y1="40" x2="15" y2="25" stroke="white" strokeWidth="2"
                                                      strokeLinecap="round"/>
                                            </svg>
                                        </button>
                                    </div>
                                )}
                            </div>
                            <div className="recipe-thumbnail-images">
                                {hasImages ? recipe.images.map((image, index) => (
                                    <div
                                        className={`recipe-thumbnail-image ${index === currentIndex ? 'active' : ''}`}
                                        key={index}
                                        onClick={() => this.handleThumbnailClick(index)}
                                    >
                                        <img src={`http://localhost:8080/${image.image_path}`}
                                             alt={`Thumbnail ${index + 1}`}/>
                                    </div>
                                )) : null}
                            </div>
                        </div>
                        <div className="recipe-details">
                            <h2>{recipe.title}</h2>
                            <p className="recipe-description">{recipe.description}</p>
                            <div className="recipe-actions-container">
                                <div className="add-mark-container">
                                    {recipe.score > 0 ? (
                                        <div className="recipe-stars-container">
                                            <h4>{(recipe.score).toFixed(1)}</h4>
                                            <span>{this.getStars(recipe.score)}</span>
                                        </div>
                                    ) : (
                                        <div className="recipe-stars-container">
                                            <h4>Нет оценок</h4>
                                        </div>
                                    )}
                                    <p className="add-mark-action" onClick={this.toggleFormVisibility}>Поставить
                                        оценку</p>

                                </div>
                                <div className="add-tried-favorite-container">
                                    <div className="heart-container" onClick={this.handleAddFavorite}>
                                        <svg
                                            className={`heart-icon ${isFavorite ? 'filled' : ''}`}
                                            viewBox="0 0 24 24"
                                            xmlns="http://www.w3.org/2000/svg"
                                        >
                                            <path
                                                d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"
                                            />
                                        </svg>
                                    </div>
                                    <div className="add-tried-container">
                                        <button className="add-tried-btn" onClick={this.handleAddTried}>
                                            {isTried ? (
                                                <div>
                                                    Опробовано
                                                </div>
                                            ) : (
                                                <div>
                                                    Опробовать
                                                </div>
                                            )}
                                        </button>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                    <div className="recipe-instructions-ingredients-container">
                        <div className="recipe-instructions-ingredients-card">
                            <div className="recipe-ingredients">
                                <h2>Ингредиенты</h2>
                                <ul>
                                    {ingredients.length ? (
                                        ingredients.map((ingredient, index) => (
                                            <li key={index}>
                                                <div className="ingredient-container">
                                                    <p>{ingredient.title}</p>
                                                    <p>{ingredient.quantity}</p>
                                                </div>
                                            </li>
                                        ))
                                    ) : (
                                        <div>
                                            <p>Ошибка загрузки ингредиентов</p>
                                        </div>
                                    )}
                                </ul>
                                <div className="recipe-shop-link-container">
                                    <Link to={`/recipe/` + recipe.id + `/shop`}>Рассчитать стоимость</Link>
                                </div>
                            </div>
                            <div className="recipe-instructions">
                                <h2>Инструкции</h2>
                                <ol>
                                    {instructions.length ? (
                                        instructions.map((instruction, index) => (
                                            <div className="instruction-container">
                                                <li key={index}>
                                                    <p>{instruction}</p>
                                                </li>
                                            </div>
                                        ))
                                    ) : (
                                        <div>
                                            <p>Ошибка загрузки инструкций</p>
                                        </div>
                                    )}
                                </ol>
                            </div>

                        </div>
                    </div>
                </div>
                <FooterComponent/>
                {isFormVisible ? (
                    <div>
                        <ScrollToTop/>
                        <FormMarkComponent
                            userScore={userScore}
                            handleRatingSubmit={this.handleRatingSubmit}
                            toggleFormVisibility={this.toggleFormVisibility}/>
                    </div>
                ) : null}
            </div>
        );
    }
}

export default withParams(Recipe);
