Ludum Dare 41
=============

## Theme: Combine 2 Incompatible Genres

I thought of a making a 2-game in one. And the two games would be interacting.
The first game is a classic beat the up: you must go throught the level without dying. The second game is a simulation game: you'r inside the player's body and have to take care about his health while he is fighting.

Win condition:

* reach the end of the level alive.

Lose condition:

* One vital organ break down.

### Controls
To differentiate the 2 games let's say that the beat the up require only keyboard controls and the body simulation only mouse controls.

### Beat them up
This is a classic beat them up, ennemies comming from the right and the player comming from the left. With a horizontal scrolling centered on the player.

He can

* Go left, right, jump, get down
* Hit when he's not jumping

He must be able to grab items

* food
* weapons ? I'm not sure about that because it doesn't involve the simulation
* pills

Those items may appear if the player kill an ennemi and if he break background elements (if i have enough time to implement it)

### Body simulation
The men is a cyborg so if we play the cyborg we also have manage his body. The software is a very good one and offer a medival interface to make things easier to understand.

Organs:

* Heart provide blood pressure so oxygen could go faster to the muscles. He need neither oxygen nor nutriments. The only way to heal the heart is to find pills.
* Lungs provide oxygen to the body. Lungs provide more oxygen if they are in good health.
* Muscles (legs and arms) requires a certain amount of oxygen to stay effective. If they stay to long without enough oxygen the player will die.

Ressources:

* Blood pressure (provided by heart)
* Nutriments (provided by food)
* Oxygens (provided by lungs)

Action:

* Can increase heartbeat to boost the player but traide some heart's life.
* Can increase respiration to get more oxygen but traide some lungs' life.
* Can buy upgrades for the muscles against nutriments.

The heart is an hospital on the center, lungs are mills and muscles are houses. Ressources are representing by peoples walking along the street.