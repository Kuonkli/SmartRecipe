import React, { Component } from "react";
import "../css/FormMarkComponent.css";
import DeleteImage from  "../assets/images/delete-image.png"

class FormMarkComponent extends Component {
    constructor(props) {
        super(props);
        this.state = {
            hoverRating: 0,
            selectedRating: props.userScore || 0
        };
    }

    componentDidMount() { document.body.classList.add('no-scroll'); }
    componentWillUnmount() { document.body.classList.remove('no-scroll'); }

    handleMouseEnter = (index) => {
        this.setState({ hoverRating: index });
    };

    handleMouseLeave = () => {
        this.setState({ hoverRating: 0 });
    };

    handleRatingSelect = (index) => {
        this.setState({ selectedRating: index });
    };

    handleSaveRating = () => {
        const selectedRating = this.state.selectedRating;
        this.props.handleRatingSubmit(selectedRating);
    }

    render() {
        const {  toggleFormVisibility } = this.props;
        const { hoverRating, selectedRating } = this.state;
        return (
            <div className="form-container">
                <div className="form-background" onClick={toggleFormVisibility}></div>
                <form className="add-mark-form">
                    <div className="form-data">
                        <div className="close-form-icon">
                            <span onClick={toggleFormVisibility}>&#10006;</span>
                        </div>
                        <h2>Оцените рецепт</h2>
                        <div className="stars-form-container">
                            <div className="stars-select-container">
                                {[1, 2, 3, 4, 5].map((index) => (
                                    <span
                                        key={index}
                                        className={`star ${index <= (hoverRating || selectedRating) ? 'filled' : ''}`}
                                        onMouseEnter={() => this.handleMouseEnter(index)}
                                        onMouseLeave={this.handleMouseLeave}
                                        onClick={() => this.handleRatingSelect(index)}
                                    >
                                        ★
                                    </span>
                                ))}
                            </div>
                        </div>
                        <div className="mark-form-button-container">
                            <button type="submit" className="edit-save-mark-form-btn" onClick={this.handleSaveRating}>Сохранить</button>
                            <button type="button" className="delete-mark-btn" onClick={() => {this.handleRatingSelect(0)}}>
                                <img src={DeleteImage} />
                            </button>
                        </div>
                    </div>
                </form>
            </div>
        );
    }
}

export default FormMarkComponent;
