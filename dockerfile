FROM golang:1.18-bullseye
# Create app folder 
RUN mkdir /app
# Copy our file in the host contianer to our contianer
COPY . /app
# Set /app to the go folder as workdir
WORKDIR /app/cmd
# Generate binary file from our /app
RUN go build
# Expose the port 8080
EXPOSE 8080:8080
# Run the app binarry file 
CMD ["./cmd"]
