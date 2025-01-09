import React, {Component, useEffect} from "react";
import { useParams } from "react-router-dom";
import HeaderComponent from "./HeaderComponent";
import FooterComponent from "./FooterComponent";
import "../css/NearestShop.css"

function withParams(Component) {
    return props => <Component {...props} params={useParams()} />;
}

const ScrollToTop = () => { useEffect(() => {
    window.scrollTo(0, 0); }, []);
    return null;
}

class NearestShop extends Component {
    constructor(props) {
        super(props);
        this.state = {
            "address": "улица Октября, вл10",
            "cart_price": 680.45,
            "client_coordinates": {
                "Lat": 55.753119,
                "Lon": 37.857009
            },
            "distance": 184.4027707,
            "message": "success",
            "products": [
                {
                    "title": "Куриное филе",
                    "approximate_price": 269.99
                },
                {
                    "title": "Лук репчатый",
                    "approximate_price": 179.99
                },
                {
                    "title": "Морковь",
                    "approximate_price": 79.99
                },
                {
                    "title": "Рис",
                    "approximate_price": 114.99
                },
                {
                    "title": "Специи",
                    "approximate_price": 35.49
                }
            ],
            "recipe": "Плов с курицей",
            "shop_coordinates": {
                "Lat": 55.751898,
                "Lon": 37.859201
            },
            "shop_id": 1,
            "shop_name": "Перекресток"
        }
    }

    componentDidMount() {
        const script = document.createElement("script");
        script.src = "https://api-maps.yandex.ru/2.1/?lang=ru_RU";
        script.async = true;
        script.onload = () => {
            window.ymaps.ready(this.initMap);
        };
        document.body.appendChild(script);
    }

    initMap = () => {
        const { client_coordinates, shop_coordinates } = this.state;

        const map = new window.ymaps.Map("map", {
            center: [
                (client_coordinates.Lat + shop_coordinates.Lat) / 2,
                (client_coordinates.Lon + shop_coordinates.Lon) / 2,
            ],
            controls: [],
            zoom: 15,
        });

        const clientPlacemark = new window.ymaps.Placemark(
            [client_coordinates.Lat, client_coordinates.Lon],
            {
                hintContent: "Клиент",
                balloonContent: "Клиент",
            },
            {
                preset: "islands#redDotIcon",
            }
        );

        const shopPlacemark = new window.ymaps.Placemark(
            [shop_coordinates.Lat, shop_coordinates.Lon],
            {
                hintContent: "Магазин",
                balloonContent: "Магазин",
            },
            {
                preset: "islands#blueDotIcon",
            }
        );

        map.geoObjects.add(clientPlacemark);
        map.geoObjects.add(shopPlacemark);
    };

    render() {
        return (
            <div className="nearest-shop">
                <HeaderComponent />
                <ScrollToTop />
                <div className="nearest-shop-page">
                    <div className="shop-map-info-container">
                        <div className="nearest-shop-info">
                            <h1 className="nearest-shop-headlines">Расчет стоимости в ближайшем магазине</h1>
                            <h2>{this.state.shop_name}</h2>
                            <h4>Адрес: {this.state.address}, расстояние: {this.state.distance.toFixed(0)} метра.</h4>
                        </div>
                        <div className="nearest-shop-map">
                        <div id="map"></div>
                        </div>
                    </div>
                    <div className="shopping-cart">
                        <h2>Цены на ингредиенты</h2>
                        <h3>Рецепт: {this.state.recipe}</h3>
                        <ul>
                            {this.state.products.map((ingredient, index) => (
                              <li>
                                  <div className="product-container">
                                      <p>{ingredient.title}</p>
                                      <p>{ingredient.approximate_price} руб.</p>
                                  </div>
                              </li>
                            ))}
                        </ul>
                        <h3>Примерная цена корзины: {this.state.cart_price} руб.</h3>
                    </div>
                </div>
                <FooterComponent />
            </div>
        );
    }
}

export default withParams(NearestShop);
