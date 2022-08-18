import Swiper from 'swiper';
import 'swiper/swiper-bundle.css';
const swiperPlugin = {
  install(Vue) {
    Vue.prototype.$swiperCreate = (...args) => new Swiper(...args);
  },
}
export default swiperPlugin;
