import React, { Component } from "react";
import HeaderComponent from "./HeaderComponent";
import DefaultUserImage from "../assets/images/default-user-image.png";
import RecipesGridComponent from "./RecipesGridComponent";
import axios from "axios";
import authService from "../services/authService";
import FooterComponent from "./FooterComponent";
import "../css/Profile.css";
import { Link } from "react-router-dom";

class Profile extends Component {
    constructor(props) {
        super(props);
        this.state = {
            userInfo: {},
            favoriteRecipes: [],
            triedRecipes: [],
            ratedRecipes: [],
            currentIndexes: {},
            isSmallScreen: window.innerWidth < 768,
            isEditMode: false,
            editedName: '',
            editedSurname: ''
        };
    }

    async componentDidMount() {
        window.addEventListener("resize", this.handleResize);
        await this.fetchFavoriteRecipes();
        await this.fetchTriedRecipes();
        await this.fetchRatedRecipes();
        await this.fetchUserInfo("http://localhost:8080/api/user/profile");
    }

    componentWillUnmount() {
        window.removeEventListener("resize", this.handleResize);
    }

    handleResize = () => {
        this.setState({ isSmallScreen: window.innerWidth < 768 });
    }

    fetchFavoriteRecipes = async () => {
        await this.fetchRecipes(`http://localhost:8080/api/user/favorite_recipes?limit=4`, 'favoriteRecipes');
    }

    fetchTriedRecipes = async () => {
        await this.fetchRecipes(`http://localhost:8080/api/user/tried_recipes?limit=4`, 'triedRecipes');
    }

    fetchRatedRecipes = async () => {
        await this.fetchRecipes(`http://localhost:8080/api/user/rated_recipes?limit=4`, 'ratedRecipes');
    }

    fetchUserInfo = async (url) => {
        try {
            const accessToken = localStorage.getItem('accessToken');
            const response = await axios.get(url, {
                headers: {
                    Authorization: `Bearer ${accessToken}`
                }
            });

            if (response.data) {
                this.setState({
                    userInfo: response.data.user,
                    editedName: response.data.user.name,
                    editedSurname: response.data.user.surname
                });
            } else {
                throw new Error("Invalid response data format");
            }
        } catch (error) {
            console.error("Failed to fetch userInfo:", error);
            if (error.response && error.response.status === 401) {
                await authService.refreshToken();
                return this.fetchUserInfo(url);
            }
        }
    }

    fetchRecipes = async (url, type) => {
        try {
            const accessToken = localStorage.getItem('accessToken');
            const response = await axios.get(url, {
                headers: {
                    Authorization: `Bearer ${accessToken}`
                }
            });

            if (response.status === 204) {
                return;
            }

            if (response.data && Array.isArray(response.data.recipes)) {
                const currentIndexes = {};
                response.data.recipes.forEach(recipe => {
                    currentIndexes[recipe.id] = 0;
                });

                this.setState({
                    [type]: response.data.recipes,
                    currentIndexes: { ...this.state.currentIndexes, ...currentIndexes }
                });
            } else {
                throw new Error("Invalid response data format");
            }
        } catch (error) {
            console.error("Failed to fetch recipes:", error);
            if (error.response && error.response.status === 401) {
                await authService.refreshToken();
                return this.fetchRecipes(url, type);
            }
        }
    }

    handlePrevImage = (recipeId) => {
        this.setState(prevState => {
            const currentIndexes = { ...prevState.currentIndexes };
            const recipe = this.findRecipeById(recipeId);
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
            const recipe = this.findRecipeById(recipeId);
            if (currentIndexes[recipeId] < recipe.images.length - 1) {
                currentIndexes[recipeId] += 1;
            } else {
                currentIndexes[recipeId] = 0;
            }
            return { currentIndexes };
        });
    }

    findRecipeById = (recipeId) => {
        return [...this.state.favoriteRecipes, ...this.state.triedRecipes, ...this.state.ratedRecipes].find(r => r.id === recipeId);
    }

    handleEditMode = () => {
        this.setState(prevState => ({
            isEditMode: !prevState.isEditMode
        }));
    }

    handleCancelEditing = () => {
        this.setState({
            editedName: this.state.userInfo.name,
            editedSurname: this.state.userInfo.surname,
            isEditMode: false
        });
    }

    handleChange = (event) => {
        const { name, value } = event.target;
        this.setState({ [name]: value });
    }

    handleLogOut = () => {
        authService.clearTokens();
        window.location.href = "/";
    }

    render() {
        const { isSmallScreen, favoriteRecipes, triedRecipes, ratedRecipes, currentIndexes, userInfo, isEditMode, editedName, editedSurname } = this.state;
        const recipesToShow = isSmallScreen ? 4 : 3;

        return (
            <div className="profile">
                <HeaderComponent />
                <div className="profile-card">
                    <div className="user-image-container">
                        <img src={DefaultUserImage} alt="Пользователь" />
                    </div>
                    <div className="user-information-container">
                        <p>Имя</p>
                        <input
                            className={`user-name ${isEditMode ? 'editing' : ''}`}
                            value={editedName}
                            name="editedName"
                            onChange={this.handleChange}
                            disabled={!isEditMode}
                        />
                        <p>Фамилия</p>
                        <input
                            className={`user-surname ${isEditMode ? 'editing' : ''}`}
                            value={editedSurname}
                            name="editedSurname"
                            onChange={this.handleChange}
                            disabled={!isEditMode}
                        />
                        <p>E-mail</p>
                        <input className="user-email" defaultValue={userInfo.email} disabled={true} />
                        <div className="user-profile-control-buttons">
                            {isEditMode ? (
                                <div className="save-user-data-container">
                                    <button className="save-user-data-btn">Сохранить</button>
                                    <button className="cancel-editing-btn" onClick={this.handleCancelEditing}>Отмена</button>
                                </div>
                            ) : (
                                <button className="edit-user-profile-btn" onClick={this.handleEditMode}>Редактировать</button>
                            )}
                            <button className="logout-user-profile-btn" onClick={this.handleLogOut}>Выйти</button>
                        </div>
                    </div>
                </div>
                <div className="personal-recipes-container">
                    <div className="profile-head-lines">
                        <h1>Избранное</h1>
                        <Link to={`/recipes/favorite`}>Смотреть больше</Link>
                    </div>
                    {this.state.favoriteRecipes.length ? (
                        <RecipesGridComponent
                            recipes={favoriteRecipes.slice(0, recipesToShow)}
                            currentIndexes={currentIndexes}
                            handlePrevImage={this.handlePrevImage}
                            handleNextImage={this.handleNextImage}
                        />
                    ) : (
                        <p className="default-profile-recipes-text">У вас пока нет избранных рецептов</p>
                    )}
                    <div className="profile-head-lines">
                        <h1>Опробованное</h1>
                        <Link to={`/recipes/tried`}>Смотреть больше</Link>
                    </div>
                    {this.state.triedRecipes.length ? (
                        <RecipesGridComponent
                            recipes={triedRecipes.slice(0, recipesToShow)}
                            currentIndexes={currentIndexes}
                            handlePrevImage={this.handlePrevImage}
                            handleNextImage={this.handleNextImage}
                        />
                    ) : (<p className="default-profile-recipes-text">У вас пока нет опробованных рецептов</p>)}
                    <div className="profile-head-lines">
                        <h1>Оцененное</h1>
                        <Link to={`/recipes/rated`}>Смотреть больше</Link>
                    </div>
                    {this.state.ratedRecipes.length ? (
                        <RecipesGridComponent
                            recipes={ratedRecipes.slice(0, recipesToShow)}
                            currentIndexes={currentIndexes}
                            handlePrevImage={this.handlePrevImage}
                            handleNextImage={this.handleNextImage}
                        />
                    ) : (<p className="default-profile-recipes-text">У вас пока нет оцененных рецептов</p>)}
                </div>
                <FooterComponent />
            </div>
        );
    }
}

export default Profile;
