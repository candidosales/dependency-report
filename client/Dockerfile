FROM node:21-alpine AS compile-image

WORKDIR /opt/ng
COPY package.json ./
RUN npm install

ENV PATH="./node_modules/.bin:$PATH" 

COPY . ./
RUN ng build --configuration=production

FROM nginx:alpine
# COPY docker/nginx/default.conf /etc/nginx/conf.d/default.conf
COPY --from=compile-image /opt/ng/dist/dependency-report /usr/share/nginx/html