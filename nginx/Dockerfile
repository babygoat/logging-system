FROM nginx:1.16.1

ARG home=/home/nginx

WORKDIR $home

# Create directory for nginx temp files
RUN mkdir -p ${home}/tmp

COPY nginx.conf cert.pem privkey.pem ./

RUN chown -R nginx:nginx $home

EXPOSE 80

USER nginx

CMD nginx -c /home/nginx/nginx.conf -g "daemon off;"
