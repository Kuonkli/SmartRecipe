import React, { Component } from "react";
import { Link } from "react-router-dom";
import DefaultRecipeImage from "../assets/images/default-recipe-image.png";
import ("../css/RecipeGrid.css")

class RecipesGrid extends Component {
    getStars(rating) {
        const fullStars = Math.floor(rating);
        const halfStar = rating % 1 >= 0.5 ? "★" : "";
        const emptyStars = "☆".repeat(5 - fullStars - (halfStar ? 1 : 0));
        return "★".repeat(fullStars) + halfStar + emptyStars;
    }

    render() {
        const { recipes, currentIndexes, handlePrevImage, handleNextImage } = this.props;
        return (
            <div className="recipes-grid">
                {recipes.map((recipe) => (
                    <div className="recipe-card" key={recipe.id}>
                        {recipe.images && recipe.images.length > 0 ? (
                            recipe.images.length > 1 ? (
                                <div className="recipe-slider-container">
                                    <button className="slider-btn" onClick={() => handlePrevImage(recipe.id)}>
                                        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 50">
                                            <line x1="12" y1="10" x2="2" y2="25" stroke="#12a370" strokeWidth="2" strokeLinecap="round"/>
                                            <line x1="12" y1="40" x2="2" y2="25" stroke="#12a370" strokeWidth="2" strokeLinecap="round"/>
                                        </svg>
                                    </button>
                                    <img
                                        src={`http://localhost:8080/${recipe.images[currentIndexes[recipe.id] || 0].image_path}`}
                                        alt={recipe.title}/>
                                    <button className="slider-btn" onClick={() => handleNextImage(recipe.id)}>
                                        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 50">
                                            <line x1="8" y1="10" x2="18" y2="25" stroke="#12a370" strokeWidth="2" strokeLinecap="round"/>
                                            <line x1="8" y1="40" x2="18" y2="25" stroke="#12a370" strokeWidth="2" strokeLinecap="round"/>
                                        </svg>
                                    </button>
                                </div>
                            ) : (
                                <img src={`http://localhost:8080/${recipe.images[0].image_path}`} alt={recipe.title}/>
                            )
                        ) : (
                            <img src={DefaultRecipeImage} alt={recipe.title}/>
                        )}
                        <h3>{recipe.title}</h3>
                            { recipe.score > 0 ? (
                                <div className="stars">
                                    <h4>{recipe.score.toFixed(1)}</h4>
                                    <span>{this.getStars(recipe.score)}</span>
                                </div>
                            ) : (
                                <div className="stars">
                                    <h4>Нет оценок</h4>
                                </div>
                            )}
                        <Link to={`/recipe/${recipe.id}`} className="see-recipe-btn">Посмотреть</Link>
                    </div>
                ))}
            </div>
        )
    }
}

export default RecipesGrid;
