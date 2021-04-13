import React from 'react';
// import sections
import Hero from '../components/sections/Hero';
import FeaturesTiles from '../components/sections/FeaturesTiles';
import FeaturesGithubSplit from '../components/sections/FeaturesGithubSplit';
import FeaturesGitlabSplit from '../components/sections/FeaturesGitlabSplit';
import Cta from '../components/sections/Cta';

const Home = () => {

  return (
    <>
      <Hero className="illustration-section-01" />
      <FeaturesTiles />
      <FeaturesGithubSplit invertMobile topDivider imageFill className="illustration-section-01" />
      <FeaturesGitlabSplit invertMobile topDivider imageFill className="illustration-section-02" />
      <Cta />
    </>
  );
}

export default Home;