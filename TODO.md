## GENERAL
- [ ] No Javascript (or optional config)
- [ ] Better Error Pages
- [ ] Search bar
- [ ] Search via flair (optional config to save to file?) (for example `subreddit-flair.txt` (or custom) and it contains flair names)

## TYPES
- [ ] Implement CrossPost
- [ ] Implement GeoFilter

## CSS
- [ ] Optimize it for smaller screens (currently best used at >=1920x1080)
- [ ] Add different styles (eg, old, solarized, or barebones)
- [X] Fix subreddit description overflowing
- [X] Fix source buttons causing overflow

## SUBREDDITS / TEMPLATE
- [ ] Implement Viewing Comments
- [ ] Implement Share button (Permalink to reddit post)
- [ ] Add `Resolution Quality` config (1->5->source)
- [ ] Implement settings page
- [X] Use SVG instead of Google Icons
- [X] Remove iframes (add a thumbnail?) (top priority)
- [ ] Source image for galleries
- [X] Limit img width to 512px(?)
- [X] (Video) images sometimes only have reddit_video_preview, and not mp4.
    - Funnily enough, even reddit made this mistake, mistaking a video for an image.
- [X] Trying to open an image in a new window on galleries only opens the *last* image of the gallery
- [ ] Comments: Achievements
- [ ] Comments: Starred Comments

## GO
- [X] Fix the horrible for loop, or just make it readable
- [X] Don't use range in MediaMetaData (or try using sort?)
- [ ] Use at least the 3rd resolution? (anything less is pretty blurry, or too small.)
- [ ] Galleries: if it's a RedditVideo, then grab the id, and use `https://v.redd.it/id/DASH_q.mp4`
- [X] Serve posts.html instead of loadedposts.html
- [X] Error Handling: Skip post if it errors out
- [X] Use MP4 instead of GIF (top priority)
    - For example, a 2 second, 1820x2000 GIF is around 7,252,000 bytes (7.2MB). Comparing it to the same width x height MP4, the result is around 587,000 bytes (587KB) (91% decrease in file size, for the same quality!)

## WORKFLOWS
- [X] Automatically append +1 to the VERSION files?