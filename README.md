# raven
**Raven** (just like the raven messengers on Westeros in Game of Throne), which is a distributed news and social networking service written in Golang. The front-end is implemented with HTML and Go, while the backend, known as [citadel](https://github.com/googlr/raven/tree/master/citadel), are multiple servers that replicate changes over the cluster to maintain consistency and availability in unreliable networks.

## Team 
  - Captain: TBD
  - Sailer 1: [xg626](https://github.com/googlr) 
  - Sailer 2: TBD
  
## Description
  The application is based on three parts:
  1. Front-end: 
  2. Back-end: `Citadel`
  
## Questions
  - Password parsing: `<input type="password" name="userpswd" value="qwerty"><br />`. I changed the `name="password"` to `name="userpswd"`, then `r.FormValue("userpswd")` is not longer empty. To be answered.
  - [iterating-through-map-in-template](https://stackoverflow.com/questions/21302520/iterating-through-map-in-template)
  
## Reference
  1. Session: [github.com/gorilla/sessions](http://www.gorillatoolkit.org/pkg/sessions)
