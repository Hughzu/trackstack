- Infra deep dive

    - understand what's done
    - how does it perform in a platform engineer scenario
    - better module structure ? 
    - better onboarding instead of bashes script ? ( a newbie has to make it work from my base framework)
    - rewrite infra.md ? 
        - need to explain deeply the sigv4 thingy because other people deploying some stuff will have a non working app. I have to document how this works quite nicely.

    - make a cli (???) to understand what the fuck is going on in bootstrap 
        - idea is to have an AWS cli set up (with root access)
            - github cli set up
            - turso cli set up
            - aws cli set up ? 
                - do I use an account for this ? right now I log in with root access on my aws console. 
        - cli have all the things at the right place in order to create a db, get back tokens and url, store them in env variable in github and / or in a specific .env file in order to initialize ssm ? 
        - then create mine that uses them in order to setup everything how I want (maybe create another AWS account for the rest, init the pipelines for terraform and things like that ?)
    
    - modify agents.md with docs about how many clis is available + the doc for my cli

- Write a blog 
    - squash commit to make this repo public ? 
    - blog focus point is have a performant 0$ infra for your 900$ tokens expense tracker app. 

- ECS Fargate

- Blog ? 
    - Sell platform engineering profile ? 

- SRL Module ? (vacancy days, checklist per month with reminders, ...)

- k3s

- EKS

- Observability / platform engineering stuff ? 