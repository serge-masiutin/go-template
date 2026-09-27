# Intent log

What we worked on, day by day, and what shipped. Short on purpose: one short
list a day per person, however many PRs came out of it.

Written by each of us with Claude, from our own session prompts. Composed from
`docs/intent/`, so edit your own file there rather than this one.

---

## Could use a hand

- **Irina**, Aug 18: an ambient background sound for the pier. I've tried twice
  and nothing feels right.
- **Irina**, Aug 14: that Slack DM deep link. Still want it, still not built.

## Wed Jul 15-16, 2026

First look. Refresh last year's product for this year's conference, sfruby.com
for context, we're on an 8-bit style now and the characters should be generated
in it.

## Sun Jul 26

The real brief. Relaunch the app for SF Ruby Startup Conference 2026, 8-bit for
the UI and the generated images both. Archive last year's ticket holders
without deleting anything. For people who came last year, reuse the photo they
already gave us so their card is ready before they show up. The near-black
background on the card pages had to go, it looked off against our design
system.

## Wed Jul 29

One app, two conferences. `#1`

## Thu Jul 30

Wrote the character prompt myself: 16-bit chibi sprite off a style reference,
tuned so people look good, because nobody shares a card that flatters them
badly. The cards can vary but the "Join..." wording is a composite overlay for
sure, we paint it, we don't ask a model for it. Rails security bump lands
before anything else. `#2`

## Fri Jul 31

Two threads. On the cards: I didn't want people surrounded by money symbols,
the coins and tickets the model added read as cash, so rubies, hearts and our
martians instead, and the background moved from clouds to any San Francisco
scene. `#3` Returning attendees get a card generated from last year's photo
automatically, before we ever write to them `#4`, and upload screening goes off
for this crowd, they're ticket holders `#5`. Vova reviewed and said the
generation code doesn't belong in `app/models/card`, it belongs in
`app/agents/<event>`, and I had it refactored before building anything else. In
a separate worktree I laid out where this is going: someone buys a ticket on
Luma, gets an invitation from us straight away, makes their card, then keeps
getting tasks and points all the way to the conference, with a gift for whoever
scores highest. Last year I exported and uploaded a CSV by hand every single
day, so I asked for real Luma integration, we're on Luma Plus, 82 people
already hold tickets, plus a profile page where people edit everything except
their email. `#8` Tried getting the card text into a shaped badge instead of a
rectangular bar, Steam capsule art as my reference, never landed it
`#7 dropped`. Brought the review tooling in, the pr-review skill from Solaris
and palkan's layered-rails skills, so reviews here are held to his conventions
and not mine.

## Tue Aug 4 - Wed Aug 5

The logo moves from the middle of the card's strip to the top left corner. Last
year's gallery is preserved on its own public page instead of vanishing, linked
from the footer, and we reuse the sfruby.com footer with the conference links.
Rewrote the gallery headline to sell the conference rather than just be cute.

## Mon Aug 10

Pre-launch. Vova reviewed the Luma work and said enrolling each guest should be
its own small job, so that went in first. Merge everything, keep 2026 switched
off, import the real ticket holders, generate their cards and review them in
admin against a checkerboard so I can see the cutouts before anyone else does
`#9` `#15`. The game economy: coins and hearts, every accrual and spend
recorded, quests for publishing and sharing `#10`. The wallet page was 500ing,
and admin review wanted full-size previews `#12`. The whole thing should feel
like one app: sign-in without passwords, the same nav and footer everywhere,
your own card as your home page, conference links in the footer `#14`. People
invited before go-live keep working links `#16`. Refused the word "veteran" for
returning attendees, we have actual veterans in the community.

*From review: `#11` `#13` `#17`*

## Tue Aug 11

Went live at about eight in the evening. `#6` Before that: the team said the
page can't look the same to you and to a stranger, so the two got split, and
the gallery became one full-size card per row with details in a modal `#19`. A
public API so sfruby.com can show the newest cards live `#18`. What we share is
sfruby.com and not us, so the utm tail comes off `#21`. Name, company, role and
intro go on the card, and the speaker and organizer tags come back off, the
cards looked busy. Sharing to X, LinkedIn and Bluesky, dropping the hashtag in
favour of mentioning the conference account. Sharing buttons and a person's own
social links stop being the same thing, and the whole thing is renamed Pixel
Card `#20`. Rewrote the invitation email myself, it greeted people twice and
explained itself, so it's one line about the first quest now. Then a scare: too
many notifications went out and I thought we'd invited everyone on Luma instead
of only ticket holders. I asked to pause the sync `#22 dropped`, checked the
numbers, found most were last year's, and closed it. Added Plausible `#24`.

*From review: `#23`*

## Wed Aug 12

Quiet day. Asked for imgproxy on the gallery images. Late at night I dropped an
isometric painting of San Francisco into the repo, and that started the pier.

## Thu Aug 13

I want our pixel people standing on that city: first you find yourself and get
10 coins, then you go looking for other people. Use the characters we already
have on production, don't generate new ones. Placement fought me all day,
characters on roofs, on walls, on top of trees, whole lawns empty, so I gave up
on working it out and drew the mask of where feet may go by hand. Separately: a
lot of characters on production have their legs cropped, and I want a pipeline
that draws the missing legs back without touching the face or anything already
painted.

## Fri Aug 14

The pier ships. `#26` Changed my mind about placement: not centre-outwards,
everyone gets a random free spot, so you actually have to search for your
friends. I want to edit the mask myself and re-run. Sponsor banners stand on
rooftops, sized by tier, Pickaxe biggest, never overlapping each other or
people. 450 cloned attendees, so I can see a full pier before the real one
fills. The city should feel alive, so martians in different outfits, a phoenix,
rubies, and a blimp flying over, and I prepared the artwork myself so placement
is ours and not a model's. Hover bubbles: sponsors say something their
marketing team would like, animals say something funny. The leg repair ships as
something I approve in admin rather than something that just happens `#25`,
with free recuts for last year's cutouts `#27`. Asked what the pier costs to
load, which is how images ended up behind imgproxy `#29`. Zoomed the map 2.5x
so the city is something you scroll into and discover `#30`, and gave it a
short address to send people to `#31`. A request from the team: when I change
my name on my profile it should change on the cards I already have, and that's
just repainting, no model involved `#32`. Wanted a "send hi" button on the card
modal that opens a Slack DM with that person, decided it's a whole feature and
postponed it.

*From review: `#28` `#36`*

*Could use a hand: that Slack DM deep link. Still want it, still not built.*

## Sat Aug 15

I don't want two fields on a person that both mean their name, so `name` and
`full_name` merge into one `#34`. The blimp is full size and flies over
everything, a small one doesn't read as real `#33`. The pier gets laid out like
Google Maps, full screen with the panels collapsed on top of it `#35`. Asked
what our links look like when pasted in Slack, which turned into og tags on the
pier `#37` and on the cards `#40`. On mobile there's no hover, so tapping a
sponsor or an animal shows the bubble first and opens the card on the second
tap `#38`, modals that overflow get fixed `#39` `#41`, and the city can be
dragged `#42`. In the evening, the big one: people can move their characters.
I'd prepared three masks for it, where you can walk, where you walk hidden
behind buildings and need an x-ray view to be seen, and a few swim routes to
the ferries, with positions streamed live over AnyCable whispers for up to 450
people. Turned down the easy version, an arrow pointing you at your target, I
wanted the hard one. People and creatures greet you automatically when you come
near, and ghost paths show only while you're walking. `#43` A player told us
the game wants music, so I asked what we can use legally and free and decided:
action sounds from a library, each animal makes its own animal's sound, quests
chime, no background music yet.

## Sun Aug 16

Walking on a phone: you press on your character and swipe the direction you
want them to go, lead and speed 300, dropped to 60 on phones `#46` `#47`, and
your own greeting bubble goes quiet while you're the one walking. Quests become
one list shown in both places, newest on top `#44`, and the quest-complete
modal gets tidied, the +1 sat too far from its heart and the button's rounded
corners broke our design system `#45`. Postponed AnyCable streams history. The
area you can grab to walk is too small on a tablet, so make it a full circle
1.5x the character's height. And I want a quest board in admin, every
participant who accepted with their quests and their balances, so I can see how
the game is actually going.

## Mon Aug 17

The quest board ships `#49`. Our invitation emails were still in last year's
plain style when we'd already built the pixel design, so they got dressed
properly `#50`. Tweaked the logo myself because SAN FRANCISCO in white
disappeared on light backgrounds, which meant refreshing the map and both
walking masks around it and re-running the spots, the ground had changed `#48`.
sfruby.com needed company and job title from our API `#53`. Then sounds: they
fire on hover, the sea lions and whale get a water splash, the parrots need
their true sound, the pelicans had something that isn't a pelican, the human
"hey" was annoying so I asked for a neutral activation sound instead, and we
credit Sonic Pi because it's a Ruby tool and that's good for us `#52`. The
ocean I could not get right at any volume `#55`, so I pulled it out entirely
`#56`. And chat, which I thought was the only gameplay feature still missing:
you walk up to somebody, click them, a "say" button, 140 characters, emoji
allowed, and one simple inbox per person showing who said what in order, both
sides of it, over AnyCable `#54`.

*From review: `#51`*

## Tue Aug 18

Still no ambient San Francisco background sound, so a rare foghorn instead,
about once a minute and a bit random `#57`. A leaderboard ranked by coins,
tucked into the pier footer `#58`. The sound button needed two presses before
you heard anything, so it says what is true now `#59`. Every heart and coin
emoji becomes our 8-bit sprites, everywhere `#60`. Came back to the chat and
asked how private it is: we should never see these in admin or anywhere else,
so they're encrypted and there's no organizer screen. Not anonymous though: we
trust attendees and every line carries a name. Positions on production drift
when they're whispered `#61`, and dragging my own character on mobile still
scrolls the map instead, so a gesture gets judged once and a thumb gets a stick
`#63`. Then the evening: the pier's JS split into a stylesheet, a geometry
object and the greeting on its own `#65`; the camera really keeping the
character in the middle of the screen `#66`; what the browser writes surviving
a render `#67`; too many menus, so quests live on the pier alone and the map
becomes the page you land on `#68`; a quest for talking to five people,
messages both ways `#69`; and what somebody said landing in the chatbox rather
than on the form `#64`. And I started keeping this log `#62`.

*Could use a hand: an ambient background sound for the pier. I've tried twice
and nothing feels right.*

## Wed Aug 19

- Svyat got a 422 clicking on his own card on production `#70`
  **why:** he had joined twice, two accounts on two different emails
- build the production image only when something it is made of changed `#71`
- share a link to the game, short, with the card as its social preview `#72`
- one address for a person, and the two it replaces redirect to it `#73`
  **I decided:** the canonical URL draws the pier behind the card, so somebody
  following a shared link lands on the game
- the card frame kept taking the page's own head with it `#74`
- og tags the platforms actually read `#75`
  **why:** the opengraph inspector called the title 88 characters and the
  description 46
- the number of talks on Rosa's rubyevents line, which just said "talks on
  rubyevents"
- reshaped the intent log on Vova's feedback: fewer words, no flourish
  **I decided:** one list of what I worked on rather than wanted and done split
  in two, and the PR status says where each line got to

## Thu Aug 20

- get a player out of where the map has them standing
- the tab title carries the name, the conference and the dates `#76`
- one reach check for the greeting and the reply, settled where the bubble
  opened `#77`
  **why:** the bubble opened on whispered positions and the send was judged on
  stored ones, so a reply typed into an open bubble came back refused
- only the greeting you can answer speaks up `#78`
  **I decided:** the animals and the martians keep their bubbles exactly as
  they are
- a tab left idle for half an hour comes back to a city that walked on `#79`
  **I decided:** redraw after five minutes away, and skip it while the reply
  box still holds words

## Mon Aug 24

- zoom the game in 2x on a phone, and take the pinch off the browser `#80`
  **why:** the browser's own zoom scales the navigation and the text with it
- the log records the human decisions and the reasoning under them, on Svyat's
  ask
  **why:** he reviews a PR and cannot tell the human's premise from the model's
  additions
  **I decided:** only human decisions go in; the model's own stay out, and
  nothing needs to say so

## Tue Aug 25

- a gem hunt: five real Ruby gem names dropped on the map a day, a coin each
  `#81`
  **I decided:** one person per gem, and a gem you maintain pays ten rather
  than one
- gems maintained by the confirmed speakers, over 100 stars
- close the message bubble once the message has gone `#82`
- returning attendees who used a different email last year `#83`
  **I decided:** match on last year's full name or the email, either one
- TRMNL as an emerald sponsor, and Fin's new logo `#84`
- a self-review mode for the review skill: address every finding, post none of
  them

## Wed Aug 26

- the fork bomb that OOM'ed my laptop a handful of times `#85`
- let somebody walled in by a repaint walk out of it `#86`
  **why:** a player wrote in that they were stuck in the sand and the arrow
  keys did nothing

## Fri Aug 28

- wrote down how to reproduce the fork bomb, since Vova could not

## Sat Aug 29

- the queue database was throwing busy timeouts all day `#87`
  **why:** Vova spotted it; the queue was keeping every guest it had ever
  enrolled
- two activation emails over the campaign and no more: a first chat message,
  and the gem hunt opening `#88`
- call the game Clouds everywhere a player can read it `#89`
  **I decided:** "See you in the Clouds" is the closing line, not the Pier and
  not the map
- grow the gem roster from the GitHub accounts the room arrives with `#92`
  **why:** the curated list runs out four weeks before the doors open
  **I decided:** rank on downloads rather than stars, since stars miss a gem
  whose repository is not where its project's stars landed

## Sun Aug 30

- the Collect button on a gem could not be clicked `#90`
- browser-side Sentry, without posting the keys `#91`
- every gem in reach offers itself `#93`
  **why:** standing beside one person left two gems at my feet with no button
  at all
- build the front end before the suite rather than inside it `#94`
- put a bubble away with a tap on a phone `#95`
  **why:** it sometimes covers something important, like myself
- share on a phone the way the desktop does `#96`
  **I decided:** X or copy link and never the picture; LinkedIn opens its own
  feed on a phone, so there is no button for it

## Mon Aug 31

- the quest rows had stopped agreeing: "10 Banked" on one, "+10" on another,
  the same thing `#97`
  **I decided:** unify on "+10" and a coin, extracted into a presenter rather
  than patched view by view
- give the outbound calls one way to fail `#98`
- finish the stranded local-only work, where the walk-out fix was missing its
  server half `#115`

## Tue Sep 1

- the pier's own object for what one request sees, so the controller answers
  for HTTP and nothing else `#99`
- move somebody the walking cannot `#100`
- the gem hunt announcement said it had worked and no email arrived `#101`
  **I decided:** add a canary, so mail nothing is sending gets noticed rather
  than sitting there
- drop the gems three times a day, and each name three times over the run
  `#102`
  **I decided:** lower the attendee gem bar from 100k downloads to 50k
- Q, F and C fold the corner panels, and space to jump `#103`
- tell the walker about the keys the moment they start working `#107`
- archspec, herb, and the updated layered-rails skills `#104`
- one invitation path out of the two that had drifted `#105`
- in admin, invite again whoever has not opened theirs or is not in the game
  `#106`
- stop counting the gems in the announcement email `#108`
- twenty gems in tonight's drop, tonight only `#109`
- leave a gem lying where it fell `#110`
  **I decided:** do not retire them at all, they are going to be hunted
- a reader for the JavaScript and the CSS `#111`
  **why:** two thirds of this repo had no linter looking at it
- three gems a person a day `#112`
  **why:** one player took the whole floor
- drop the gems in the hardest places on the map to walk to `#113`
- the talk quest and the gem names on the admin quests dashboard `#114`
  **why:** Vladimir showed fewer coins than me having done more

## Wed Sep 2

- the two rules left over from the refactoring plan, and the hooks wired to
  nothing `#116`
- Order.co on a downtown roof `#117`
- watch the tooling, the tech debt and the regressions `#118`
  **I decided:** Sentry watches the running app and a rake task watches the
  repository; they cannot do each other's job

## Thu Sep 3

- stop policing how far a step came `#119`
  **why:** a deploy left people further out than one step allowed and every
  write after it was refused; Vova, reviewing the walking, said just take
  what comes, who are we defending against
- trim every comment in the repo to five lines, and a linter to pin it `#120`
- drop `pier_moved_at`, once the step check had deployed `#121`
- gift a gem you collected to somebody else, off their chat bubble `#122`
  **I decided:** the coins recalculate as if the receiver had collected it, so
  the gifter gives back what they gained
  **I decided:** the daily allowance does not move with a gift, unlike the
  coins
  **I decided:** a gift lands in the receiver's chatbox, or they would never
  hear of it
  **I decided:** losing the race to a gem chimes like winning it, and that
  is fine
- from testing on my phone: a gem's Collect button took no click, and your
  own gem played a sad chord
  **I decided:** `quest-complete.mp3` goes altogether, the usual gem sound
  for your own, and a sound for gifting too
- a message when the day's allowance is spent
- asked for a refactor pass over the gifting code
- everybody's ghost and swim filters, not only my own `#123`
- un-collapse the gem list on a leaderboard row, and stop truncating names
  on a phone with room to spare `#126`
  **why:** I gifted action_policy to Vladimir on prod and neither his row nor
  his coins showed it
- run the linters locally, the same ones CI runs
- caught this log up, and fixed the skill upstream where the sub-lines had
  never been published `#124`
  **why:** it stopped after Aug 18 while the hook went on staging every day
- made the log work for a team, and moved this repo onto it
  **I decided:** a Stop hook that nags when the log falls behind
- name the year in a heading, so a log read next January is not redated
  **why:** nothing wrote one down, so both scripts fell back to the year you
  happen to read it in, and every weekday in it goes wrong
  **I decided:** a file each, composed into the one everybody reads, so two
  people writing the same day never conflict

## Fri Sep 4

- on the same branch: a return-a-gem quest, a GitHub handle quest, five gems
  six times a day
  **I decided:** walking a gem home pays the carrier ten as well as the
  author, for every gem walked home, once per gem name rather than per row,
  so the same author cannot be handed three anycables for thirty coins
  **I decided:** the handle is its own quest, worth ten, rather than a fifth
  field on the profile one; and five goes per gem name
- High Scores is the Leaderboard, opened over the city with a Ruby Goodies
  tab, keeping its address `#127`
  **why:** an easy way back, and two corner panels open at once sat on top
  of each other
  **I decided:** the goodies tab is empty for now: coins will buy items,
  physical or digital, from Ruby companies, and it says where to send yours
- the short-drop warning reaches Sentry, and the roster warning gives five
  days' notice, in days rather than drops
  **why:** three days is what curating names by hand takes, so five

## Sat Sep 5

- a gem's Collect button that did nothing on production `#128`
  **why:** I stood right next to one and pressed thirty times
- checked how I had hit three gems a day
  **why:** two of them were at half past one in the morning, and the day turns
  at five in the afternoon Pacific
- where the sign-up flow loses people, and the fixes `#129`
  **why:** 122 invited, 74 with a card, 53 in the gallery, and people have to
  get into the game
  **I decided:** the returning-attendee letter, the silent upload failure,
  clicks in the webhook and landing on the map after publishing, in that order

## Sun Sep 6

- HEIC in the file picker, once the image's libvips was shown to read one
- caught this log up
