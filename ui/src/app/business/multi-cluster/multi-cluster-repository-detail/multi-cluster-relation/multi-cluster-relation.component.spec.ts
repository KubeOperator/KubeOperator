import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MultiClusterRelationComponent } from './multi-cluster-relation.component';

describe('MultiClusterRelationComponent', () => {
  let component: MultiClusterRelationComponent;
  let fixture: ComponentFixture<MultiClusterRelationComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ MultiClusterRelationComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(MultiClusterRelationComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
