import React, { Component } from "react";
import HeaderComponent from "./HeaderComponent";
import FooterComponent from "./FooterComponent";
import axios from "axios";
import authService from "../services/authService";
import RecipesGridComponent from "./RecipesGridComponent";
import "../css/Recipes.css"
import SortOptionIcon from "../assets/images/sort-icon.png"

class Recipes extends Component {
    constructor(props) {
        super(props);
        this.state = {
            recipes: [],
            searchTerm: "",
            currentIndexes: {},
            limit: 1,
            hasMoreRecipes: true,
            sortOption: "score",
            isDropdownOpen: false
        };
    }

    async componentDidMount() {
        await this.fetchRecipes(`http://localhost:8080/api/recipes/list?limit=${this.state.limit*12}&filter=${this.state.sortOption}`);
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
                    hasMoreRecipes: response.data.has_more  // Обновляем состояние
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

    loadMoreRecipes = async () => {
        const newLimit = this.state.limit + 1;
        const searchTerm = this.state.searchTerm;
        if (searchTerm !== "") {
            await this.fetchRecipes(`http://localhost:8080/api/recipes/list?limit=${newLimit*12}&search=${searchTerm}`);
        }
        else {
            await this.fetchRecipes(`http://localhost:8080/api/recipes/list?limit=${newLimit*12}&filter=${this.state.sortOption}`);
        }
        this.setState({ limit: newLimit });
    }

    handleSearchInputChange = (event) => {
        this.setState({
            searchTerm: event.target.value
        });
    }

    handleSearchSubmit = async () => {
        const searchTerm = this.state.searchTerm;
        this.setState({
            recipes: [],
            limit: 1,
            hasMoreRecipes: true
        });
        if (searchTerm !== "") {
            await this.fetchRecipes(`http://localhost:8080/api/recipes/list?limit=12&search=${searchTerm}&filter=${this.state.sortOption}`);
        }
        else {
            await this.fetchRecipes(`http://localhost:8080/api/recipes/list?limit=12&filter=${this.state.sortOption}`);
        }
    }

    handleSortChange = async (sortOption) => {
        this.setState({
            sortOption,
            recipes: [],
            limit: 1,
            hasMoreRecipes: true,
            isDropdownOpen: false
        });
        await this.fetchRecipes(`http://localhost:8080/api/recipes/list?limit=12&search=${this.state.searchTerm}&filter=${sortOption}`);
    }

    toggleDropdown = () => {
        this.setState(prevState => ({
            isDropdownOpen: !prevState.isDropdownOpen
        }));
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
        const sortOptions = [
            { value: 'score', label: 'По рейтингу' },
            { value: 'recent', label: 'Новое' },
            { value: 'high-caloric', label: 'Более калорийное' },
            { value: 'low-caloric', label: 'Менее калорийное' },
        ];

        return (
            <div className="recipes">
                <HeaderComponent />
                <div className="head-container">
                    <h1>Наши рецепты</h1>
                    <div className="search-sort-container">
                        <div className="search">
                            <input
                                type="text"
                                value={this.state.searchTerm}
                                onChange={this.handleSearchInputChange}
                                placeholder="Поиск рецептов"
                                className="search-bar"
                            />
                            <button onClick={this.handleSearchSubmit} className="search-button">
                                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
                                    <circle cx="10" cy="10" r="8" stroke="white" strokeWidth="2" fill="none"/>
                                    <line x1="16" y1="16" x2="20" y2="20" stroke="white" strokeWidth="2" strokeLinecap="round"/>
                                </svg>
                            </button>
                        </div>
                        <div className="sort">
                            <button className="sort-btn" onClick={this.toggleDropdown}>
                                <div className="sort-option-container">
                                    <div>
                                        <p>{sortOptions.find(option => option.value === this.state.sortOption).label}</p>
                                    </div>
                                    <div>
                                        <img src={SortOptionIcon} alt={"sort"}/>
                                    </div>
                                </div>
                            </button>
                            {this.state.isDropdownOpen && (
                                <ul className="sort-dropdown">
                                    {sortOptions.map(option => (
                                        <li key={option.value} onClick={() => this.handleSortChange(option.value)}>
                                            <div>{option.label}</div>
                                            <svg xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 24 24">
                                                <circle cx="12" cy="12" r="8" stroke="#0e8c54" strokeWidth="2"
                                                        fill={this.state.sortOption === option.value ? "#12a370" : "none"}/>
                                            </svg>
                                        </li>
                                    ))}
                                </ul>
                            )}
                        </div>
                    </div>
                </div>
                <RecipesGridComponent
                    recipes={this.state.recipes}
                    currentIndexes={this.state.currentIndexes}
                    handlePrevImage={this.handlePrevImage}
                    handleNextImage={this.handleNextImage} />
                {this.state.hasMoreRecipes && (
                    <div className="see-more-container">
                        <button className="see-more-recipes-btn" onClick={this.loadMoreRecipes}>
                            Смотреть больше
                        </button>
                    </div>
                )}
                <FooterComponent />
            </div>
        );
    }
}

export default Recipes;
