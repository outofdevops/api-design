let chai = require('chai');
let chaiHttp = require('chai-http');
chai.should();


chai.use(chaiHttp);
describe('Tweets', () => {
  
  describe('/POST tweet', () => {
    it('it should return the feed', (done) => {
      chai.request('http://localhost:3000')
          .post('/feed')
          .send({username: 'Anto', tweet: "Hi 👋 from OutOfDevOps"})
          .end((err, res) => {
            res.should.have.status(200);
            res.body.should.be.a('array');
            res.should.have.header("content-type", "application/json");
            res.body[0].should.be.a('object');
            res.body[0].should.have.property('username');
            res.body[0].should.have.property('tweet');
            done();
          });
    });
  });
});
